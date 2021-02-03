# Workspace Data Source
`powerbi_workspace` represents a workspace within Power BI (also called a Group)

## Example Usage
```hcl
data "powerbi_workspace" "myworkspace" {
  name = "Sample workspace"
}

output myworkspace_id {
  value = data.powerbi_workspace.myworkspace.id
}
```



## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required) Name of the workspace.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the workspace.
<!-- docgen:ComputedParameters -->
* `capacity_id` - (Optional) Capacity ID to be assigned to workspace.
<!-- /docgen -->
