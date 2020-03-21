# powerbi_workspace
powerbi_workspace represents a worksapce within Power BI (also called a Group)

## Example Usage
``` terraform
resource "powerbi_workspace" "myworkspace" {
	name = "Sample workspace"
}
```

## Argument Reference
The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required) Name of the workspace.
<!-- /docgen -->

## Attributes Reference
The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the workspace.
<!-- docgen:ComputedParameters -->

<!-- /docgen -->