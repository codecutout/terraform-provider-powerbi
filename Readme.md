# Terraform Provider for Power BI

The Power BI Provider supports Terraform 0.12.x. It may still function on earlier versions but has only been tested on 0.12.x

* [Terraform Website](https://www.terraform.io)

## Usage Example

```
# Configure the Power BI Provider
provider "powerbi" {
  # tenant_id       = "..."
  # client_id       = "..."
  # client_secret   = "..."
  # username        = "..."
  # password        = "..."
}

# Create a workspace
resource "powerbi_workspace" "example" {
  name     = "Example Workspace"
}

# Create a pbix within the workspace
TODO
```
## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.12.x +
* [Go](https://golang.org/doc/install) version 1.13.x (to build the provider plugin)

If you're on Windows you'll also need:
* [Git Bash for Windows](https://git-scm.com/download/win)

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To build you can run:
```sh
$ go build
...

To run tests for the provider you can run:
```sh
$ go test ./...
```

The majority of tests in the provider are Acceptance Tests - which provisions real resources in power BI. It's possible to run the acceptance tests with the above command by setting the following enviornment variables: 
- `TF_ACC=1`
- `POWERBI_TENANT_ID`
- `POWERBI_CLIENT_ID`
- `POWERBI_CLIENT_SECRET`
- `POWERBI_USERNAME`
- `POWERBI_PASSWORD`

To test the plugin with terraform
- Place `terraform.exe` in `$GOPATH/bin` (which should be in path)
- Run `go install` - This will build and deploy `terraform-provider-powerbi.exe` into `$GOPATH/bin`
- When running `terraform.exe` the Power BI provider will be available


## Authentication and Authorization

TODO