# Workspace Resource
`powerbi_workspace` represents a worksapce within Power BI (also called a Group)

## Example Usage (without dedicated capacity)
```hcl
resource "powerbi_workspace" "myworkspace" {
	name = "Sample workspace"
}
```

## Example Usage (with dedicated capacity)
```hcl
resource "powerbi_workspace" "myworkspace" {
	name = "Sample workspace"
	capacity_id = "0000-1111-1111-1111"
}
```

## Argument Reference
The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required) Name of the workspace.
* `capacity_id` - (Optional) Capacity ID, which will be assigned to workspace.
<!-- /docgen -->

## Attributes Reference
The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the workspace.
<!-- docgen:ComputedParameters -->

<!-- /docgen -->