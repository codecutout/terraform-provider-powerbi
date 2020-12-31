# Workspace Resource
`powerbi_workspace` represents a workspace within Power BI (also called a Group)

## Example Usage
```hcl
resource "powerbi_workspace" "myworkspace" {
  name = "Sample workspace"
}
```

~> Renaming a workspace will delete the old workspace and create a new workspace. Power BI APIs do not provide a way to update a workspace name. In order to maintain bookmarks and user applied configuration it is strongly recommended to perform renames manually through the UI prior to running terraform

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required, Forces new resource) Name of the workspace.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the workspace.
<!-- docgen:ComputedParameters -->

<!-- /docgen -->