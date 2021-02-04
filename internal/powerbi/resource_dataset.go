package powerbi

import (
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceDataset represents a Power BI dataset
func ResourceDataset() *schema.Resource {
	return &schema.Resource{
		Create: createDataset,
		Read:   readDataset,
		Update: updateDataset,
		Delete: deleteDataset,

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("table", func(old, new, meta interface{}) bool {
				// We can update changes to existing tables, but creating new tables
				// or deleting old tables requires forcing a new dataset to be created
				oldList := old.(*schema.Set).List()
				newList := new.(*schema.Set).List()

				// if there are different number of items then items have been added
				// or removed so we need to force new
				if len(oldList) != len(newList) {
					return true
				}

				// even if the number of items are the same need to checking they all have the
				// same names. If there are new or removed names need to force new
				intersect := intersect(newList, oldList, func(a interface{}, b interface{}) bool {
					return a.(map[string]interface{})["name"] == b.(map[string]interface{})["name"]
				})
				hasNumberOfItemsChanged := len(newList) != len(intersect)
				return hasNumberOfItemsChanged
			}),
		),

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "Workspace ID in which the dataset will be added.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Dataset.",
			},
			"default_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The dataset mode or type. Any value from `push`, `pushStreaming` or `streaming`. `asAzure` and `asOnPrem` are not supported",
				ValidateFunc: validation.StringInSlice([]string{"push", "pushStreaming", "streaming"}, false),
			},
			"default_retention_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "none",
				Description:  "The dataset mode or type. Any value from `none` or `basicFIFO`",
				ValidateFunc: validation.StringInSlice([]string{"none", "basicFIFO"}, false),
			},

			"table": {
				Type:        schema.TypeSet,
				Description: "The dataset tables. Creating new tables or removing existing tables will force a new dataset to be created",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The table name",
							Required:    true,
						},
						"column": {
							Type:        schema.TypeSet,
							Description: "The column schema for this table",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "The column name",
										Required:    true,
									},
									"data_type": {
										Type:         schema.TypeString,
										Description:  "The column data type. Any value from `int64`, `double`, `bool`, `datetime`, `string` or `decimal`.",
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"int64", "double", "bool", "datetime", "string", "decimal"}, false),
									},
									"format_string": {
										Type:        schema.TypeString,
										Description: "The format of the column as specified in [FORMAT_STRING](https://docs.microsoft.com/en-us/analysis-services/multidimensional-models/mdx/mdx-cell-properties-format-string-contents)",
										Optional:    true,
									},
								},
							},
						},
						"measure": {
							Type:        schema.TypeSet,
							Description: "The measures within this table",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "The measure name",
										Required:    true,
									},
									"expression": {
										Type:        schema.TypeString,
										Description: "The DAX expression for the measure",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"relationship": {
				Type:        schema.TypeSet,
				Description: "The dataset relationships",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The relationship name and identifier",
							Required:    true,
							ForceNew:    true,
						},
						"cross_filtering_behavior": {
							Type:         schema.TypeString,
							Description:  "The filter direction of the relationship. Any value from `automatic`, `bothDirections` or `oneDirection`",
							Optional:     true,
							ForceNew:     true,
							Default:      "automatic",
							ValidateFunc: validation.StringInSlice([]string{"automatic", "bothDirections", "oneDirection"}, false),
						},
						"from_table": {
							Type:        schema.TypeString,
							Description: "The name of the foreign key table",
							Required:    true,
							ForceNew:    true,
						},
						"from_column": {
							Type:        schema.TypeString,
							Description: "The name of the foreign key column",
							Required:    true,
							ForceNew:    true,
						},
						"to_table": {
							Type:        schema.TypeString,
							Description: "The name of the primary key table",
							Required:    true,
							ForceNew:    true,
						},
						"to_column": {
							Type:        schema.TypeString,
							Description: "The name of the primary key column",
							Required:    true,
							ForceNew:    true,
						},
					},
				},
			},
		},
	}
}

func canonicalDefaultMode(value string) string {
	// DefaultMode is the only enum that PowerBI does not treat as case insensitive
	// mapping out the case insensitive value to its canonical value
	switch strings.ToLower(value) {
	case "asazure":
		return "AsAzure"
	case "asonprem":
		return "AsOnPrem"
	case "push":
		return "Push"
	case "pushstreaming":
		return "PushStreaming"
	case "streaming":
		return "Streaming"
	default:
		return value
	}
}

func createDataset(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	groupID := d.Get("workspace_id").(string)
	defaultRetentionPolicy := d.Get("default_retention_policy").(string)

	resp, err := client.PostDatasetInGroup(groupID, defaultRetentionPolicy, powerbiapi.PostDatasetInGroupRequest{
		Name:        d.Get("name").(string),
		DefaultMode: canonicalDefaultMode(d.Get("default_mode").(string)),
		Tables: genericMap(d.Get("table").(*schema.Set).List(), func(tableValues interface{}) powerbiapi.PostDatasetInGroupRequestTable {
			tableValuesMap := tableValues.(map[string]interface{})
			return powerbiapi.PostDatasetInGroupRequestTable{
				Name: tableValuesMap["name"].(string),

				Columns: genericMap(tableValuesMap["column"].(*schema.Set).List(), func(columnValues interface{}) powerbiapi.PostDatasetInGroupRequestTableColumn {
					columnValuesMap := columnValues.(map[string]interface{})
					return powerbiapi.PostDatasetInGroupRequestTableColumn{
						Name:         columnValuesMap["name"].(string),
						DataType:     columnValuesMap["data_type"].(string),
						FormatString: columnValuesMap["format_string"].(string),
					}
				}).([]powerbiapi.PostDatasetInGroupRequestTableColumn),

				Measures: genericMap(tableValuesMap["measure"].(*schema.Set).List(), func(measureValues interface{}) powerbiapi.PostDatasetInGroupRequestTableMeasure {
					measureValuesMap := measureValues.(map[string]interface{})
					return powerbiapi.PostDatasetInGroupRequestTableMeasure{
						Name:       measureValuesMap["name"].(string),
						Expression: measureValuesMap["expression"].(string),
					}
				}).([]powerbiapi.PostDatasetInGroupRequestTableMeasure),
			}
		}).([]powerbiapi.PostDatasetInGroupRequestTable),

		Relationships: genericMap(d.Get("relationship").(*schema.Set).List(), func(relationshipValues interface{}) powerbiapi.PostDatasetInGroupRequestRelationship {
			relationshipValuesMap := relationshipValues.(map[string]interface{})
			return powerbiapi.PostDatasetInGroupRequestRelationship{
				Name:                   relationshipValuesMap["name"].(string),
				FromTable:              relationshipValuesMap["from_table"].(string),
				FromColumn:             relationshipValuesMap["from_column"].(string),
				ToTable:                relationshipValuesMap["to_table"].(string),
				ToColumn:               relationshipValuesMap["to_column"].(string),
				CrossFilteringBehavior: relationshipValuesMap["cross_filtering_behavior"].(string),
			}
		}).([]powerbiapi.PostDatasetInGroupRequestRelationship),
	})
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return readDataset(d, meta)
}

func readDataset(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	groupID := d.Get("workspace_id").(string)

	dataset, err := client.GetDatasetInGroup(groupID, d.Id())
	if isHTTP404Error(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.SetId(dataset.ID)
	d.Set("name", dataset.Name)

	return nil
}

func updateDataset(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("table") {
		client := meta.(*powerbiapi.Client)

		groupID := d.Get("workspace_id").(string)
		datasetID := d.Id()

		old, new := d.GetChange("table")

		// we can only update tables that already existed
		// anything outside of these should have forced a new resource
		changedSet := new.(*schema.Set).Difference(old.(*schema.Set))
		tablesToUpdate := intersect(changedSet.List(), old.(*schema.Set).List(), func(a interface{}, b interface{}) bool {
			return a.(map[string]interface{})["name"] == b.(map[string]interface{})["name"]
		})

		for _, tableToUpdateObj := range tablesToUpdate {
			tableToUpdate := tableToUpdateObj.(map[string]interface{})
			err := client.PutTableInGroup(groupID, datasetID, tableToUpdate["name"].(string), powerbiapi.PutTableInGroupRequest{

				Name: tableToUpdate["name"].(string),

				Columns: genericMap(tableToUpdate["column"].(*schema.Set).List(), func(columnValues interface{}) powerbiapi.PutTableInGroupRequestTableColumn {
					columnValuesMap := columnValues.(map[string]interface{})
					return powerbiapi.PutTableInGroupRequestTableColumn{
						Name:         columnValuesMap["name"].(string),
						DataType:     columnValuesMap["data_type"].(string),
						FormatString: columnValuesMap["format_string"].(string),
					}
				}).([]powerbiapi.PutTableInGroupRequestTableColumn),

				Measures: genericMap(tableToUpdate["measure"].(*schema.Set).List(), func(measureValues interface{}) powerbiapi.PutTableInGroupRequestTableMeasure {
					measureValuesMap := measureValues.(map[string]interface{})
					return powerbiapi.PutTableInGroupRequestTableMeasure{
						Name:       measureValuesMap["name"].(string),
						Expression: measureValuesMap["expression"].(string),
					}
				}).([]powerbiapi.PutTableInGroupRequestTableMeasure),
			})

			if err != nil {
				return err
			}
		}
	}

	return readDataset(d, meta)
}

func deleteDataset(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	groupID := d.Get("workspace_id").(string)
	return client.DeleteDatasetInGroup(groupID, d.Id())
}
