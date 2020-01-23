package powerbi

import (
	//"crypto/md5"
	"encoding/base64"
	"github.com/alex-davies/terraform-provider-powerbi/powerbi/internal/api"
	"github.com/hashicorp/terraform/helper/schema"
	"io"
	"strings"
	"time"
)

// ResourcePBIX represents a Power BI PBIX file
func ResourcePBIX() *schema.Resource {
	return &schema.Resource{
		Create: createPBIX,
		Read:   readPBIX,
		Update: updatePBIX,
		Delete: deletePBIX,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Workspace in which the PBIX will be added",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the pbix. This will be used as the name for the report and dataset",
			},
			"content_base64": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pbix file as a base64 string",
				// StateFunc: func(val interface{}) string {
				// 	contentBytes, _ := base64.StdEncoding.DecodeString(val.(string))
				// 	hashBytes := md5.Sum(contentBytes)
				// 	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes[:])
				// 	return "md5:" + hashBase64

				// },
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Seconds to wait while publishing the pbix",
				Default:     30,
			},
			"report_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID for the report that was deployed as part of the pbix",
			},
			"dataset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID for the dataset that was deployed as part of the pbix",
			},
		},
	}
}

func getContentReader(d *schema.ResourceData) io.Reader {
	contentBase64 := d.Get("content_base64").(string)
	contentReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(contentBase64))

	return contentReader
}

func createPBIX(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	resp, err := client.PostImportInGroup(api.PostImportInGroupRequest{
		GroupID:            d.Get("workspace").(string),
		DatasetDisplayName: d.Get("name").(string),
		NameConflict:       "CreateOrOverwrite",
		Data:               getContentReader(d),
	})
	if err != nil {
		return err
	}

	im, err := client.WaitForImportToSucceed(resp.ID, time.Duration(d.Get("timeout_seconds").(int))*time.Second)
	if err != nil {
		return err
	}

	d.SetId(im.ID)
	if len(im.Reports) >= 1 {
		d.Set("report_id", im.Reports[0].ID)
	}
	if len(im.Datasets) >= 1 {
		d.Set("dataset_id", im.Datasets[0].ID)
	}
	d.Set("workspace", d.Get("workspace").(string))
	d.Set("name", d.Get("name").(string))
	d.Set("content_base64", d.Get("content_base64").(string))

	return nil

}

func readPBIX(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*api.Client)

	return nil
}

func updatePBIX(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*api.Client)

	return nil
}

func deletePBIX(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	if reportID := d.Get("report_id"); reportID != nil {
		err := client.DeleteReport(api.DeleteReportRequest{
			ReportID: reportID.(string),
		})
		if err != nil {
			return err
		}
	}

	if datasetID := d.Get("dataset_id"); datasetID != nil {
		err := client.DeleteDataset(api.DeleteDatasetRequest{
			DatasetID: datasetID.(string),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
