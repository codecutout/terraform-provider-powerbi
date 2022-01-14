# Capacity Data Source
`powerbi_capacity` represents a capacity within Power BI.
Capacity could be created in three different ways:
* Reserved Premium Per user Capacity - created by Power BI service
* Power BI Embedded - Azure resource could be created via Azure Portal, `az` command line or https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/powerbi_embedded
* Power BI Premium - purchased via Microsoft 365


## Example Usage
```hcl
data "powerbi_capacity" "mycapacity" {
  name = "Sample capacity"
}

output mycapacity_id {
  value = data.powerbi_capacity.mycapacity.id
}
```



## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required) Name of the capacity.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the capacity.
<!-- docgen:ComputedParameters -->
* `region` - (Optional) Region of the capacity.
* `sku` - (Optional) SKU of the capacity.
<!-- /docgen -->
