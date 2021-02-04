# Dataset Resource
`powerbi_dataset` represents a push dataset Power BI. 

## Example Usage

### Datasource
```hcl
resource "powerbi_workspace" "test" {
  name = "Entries Workspace"
}
resource "powerbi_dataset" "test" {
  workspace_id = powerbi_workspace.test.id
  default_mode = "push"
  name = "Entries Dataset"

  table {
    name = "entries"
    column {
      name = "entryId"
      data_type = "string"
    }
  }

  table {
    name = "entries-audit"
    column {
      name = "entryId"
      data_type = "string"
    }
    column {
      name = "modifiedBy"
      data_type = "string"
    }
  }

  relationship {
    name = "entries to entires-audit"
    from_table = "entries"
    from_column = "entryId"
    to_table = "entries-audit"
    to_column = "entryId"
    cross_filtering_behavior = "automatic"
  }
}
```

~> Due to Power BI API limitations the _only_ operation that will perform an in place update is modifying existing tables. Take care if adding or removing tables or modifying any other properties as the dataset will be deleted and recreated, this will break dependant reports.

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `default_mode` - (Required, Forces new resource) The dataset mode or type. Any value from `push`, `pushStreaming` or `streaming`. `asAzure` and `asOnPrem` are not supported.
* `name` - (Required, Forces new resource) Name of the Dataset.
* `workspace_id` - (Required, Forces new resource) Workspace ID in which the dataset will be added.
* `table` - (Required) The dataset tables. Creating new tables or removing existing tables will force a new dataset to be created. A [`table`](#a-table-block-supports-the-following) block is defined below.
* `default_retention_policy` - (Optional, Default: `none`, Forces new resource) The dataset mode or type. Any value from `none` or `basicFIFO`.
* `relationship` - (Optional, Forces new resource) The dataset relationships. A [`relationship`](#a-relationship-block-supports-the-following) block is defined below.

---

#### A `table` block supports the following:
* `name` - (Required) The table name.
* `column` - (Optional) The column schema for this table. A [`column`](#a-column-block-supports-the-following) block is defined below.
* `measure` - (Optional) The measures within this table. A [`measure`](#a-measure-block-supports-the-following) block is defined below.

---

#### A `column` block supports the following:
* `data_type` - (Required) The column data type. Any value from `int64`, `double`, `bool`, `datetime`, `string` or `decimal`.
* `name` - (Required) The column name.
* `format_string` - (Optional) The format of the column as specified in [FORMAT_STRING](https://docs.microsoft.com/en-us/analysis-services/multidimensional-models/mdx/mdx-cell-properties-format-string-contents).

---

#### A `measure` block supports the following:
* `expression` - (Required) The DAX expression for the measure.
* `name` - (Required) The measure name.

---

#### A `relationship` block supports the following:
* `from_column` - (Required, Forces new resource) The name of the foreign key column.
* `from_table` - (Required, Forces new resource) The name of the foreign key table.
* `name` - (Required, Forces new resource) The relationship name and identifier.
* `to_column` - (Required, Forces new resource) The name of the primary key column.
* `to_table` - (Required, Forces new resource) The name of the primary key table.
* `cross_filtering_behavior` - (Optional, Default: `automatic`, Forces new resource) The filter direction of the relationship. Any value from `automatic`, `bothDirections` or `oneDirection`.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the dataset.
<!-- docgen:ComputedParameters -->

<!-- /docgen -->