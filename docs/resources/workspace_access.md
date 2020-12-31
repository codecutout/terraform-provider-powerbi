# Workspace Access Resource
`powerbi_workspace_access` represents workspace access for Azure users, Apps, and security groups.


## Example Usage
```hcl
resource "powerbi_workspace_access" "allow_email_address" {
  workspace_id            = "470b0d57-1f23-4332-a16f-9235bd174318"
  group_user_access_right = "Member"
  email_address           = "powerbiuser@mycompany.com"
  principal_type          = "User"
}

resource "powerbi_workspace_access" "allow_azure_app" {
  workspace_id            = "470b0d57-1f23-4332-a16f-9235bd174318"
  group_user_access_right = "Admin"
  principal_type          = "App"
  identifier              = "1f69e798-5852-4fdd-ab01-33bb14b6e934
}
```

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `workspace_id` - (Required, Forces new resource) Workspace ID to which user access would be given.
* `group_user_access_right` - (Required) User access level to workspace. Any value from `Admin`, `Contributor`, `Member`, `Viewer` or `None`.
* `principal_type` - (Required) The principal type. Any value from `App`, `Group` or `User`.
* `email_address` - (Optional, Forces new resource) Email address of the user.
<!-- /docgen -->
<!-- docgen:ComputedParameters -->
* `identifier` - (Optional, Forces new resource) Identifier of the principal.
* `display_name` - (Optional) Display name of the principal.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the allowed user access.
<!-- docgen:ComputedParameters -->
* `identifier` - (Optional, Forces new resource) Identifier of the principal.
* `display_name` - (Optional) Display name of the principal.
<!-- /docgen -->