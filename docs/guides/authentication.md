# Authentication

The Power BI terraform provider support authenticating either by a service principal or by using user credentials. Each method requires an initial setup.

## Power BI Service Principal

The Power BI terraform provider can use a service principal to create and manage resources. This can reduce the overhead of managing Power BI users and their associated credentials.

### Creating the service principal

1. Create an Azure Active Directory app.
   * Get the `tenant_id` and `client_id` from your Azure Active Directory App
   * Generate an `client_secret` for your Azure Active Directory App
2. Create an Azure Active Directory security group.
   * Assign your Azure Active Directory App service principal to the security group
3. Enable the Power BI service admin settings.
   * In the Power BI Admin Portal under Tenant Settings enable "Allow service principals to use Power BI API
   * Specify the security group created in step 2 to restrict access to the Azure Active Directory App

Detailed instructions, including screenshots, can be found at https://docs.microsoft.com/en-us/power-bi/developer/embedded/embed-service-principal


### Configure the provider

Set the Power BI terraform provider arguments `tenant_id`, `client_id` and `client_secret` to be the values retrieved from your Azure Active Directory App. Service Principal authentication does *not* require a `username` or `password`.

```hcl
provider "powerbi" {
  tenant_id     = <tenant id from app registration>
  client_id     = <client id from app registration>
  client_secret = <client secret from app registration>
}
```

## Power BI User

An alternative administrative setup is to create a Power BI user that is only intended to be used by the terraform provider. This was previously the only way to use the Power BI APIs.

### Creating a Power BI User

1. Create an Azure Active Directory App through the wizard at https://dev.powerbi.com/apps
   * Select all permissions
   * Take a note of the `client_id` and `client_secret`
2. Determine the `tenant_id` by looking at the Azure Active Directory app registration that was created in step 1 from within the Azure Portal
3. Create a Power BI user
   * Link in Power BI Admin Portal under the User section will direct you to the Office 365 page for creating new users
   * Ensure user has access to Power BI
   * Ensure MFA is disabled for the created user
   * Take note of the `username` and `password` of your newly created user

### Configure the provider

Set the Power BI terraform provider arguments `tenant_id`, `client_id`, `client_secret`, `username` and `password` to be the values retrieved from your Azure Active Directory App and user creation.

```hcl
provider "powerbi" {
  tenant_id     = <tenant id from app registration>
  client_id     = <client id from app registration>
  client_secret = <client secret from app registration>
  username      = <username from powerbi user>
  password      = <username from powerbi user>
}
```