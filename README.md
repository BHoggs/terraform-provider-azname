# Terraform Provider for Azure Resource Naming (azname)

A [Terraform](https://www.terraform.io) provider for generating standardized Azure resource names following naming conventions. This provider helps maintain consistent resource naming across your Azure infrastructure by providing configurable templates and functions for generating resource names.

## Features

- **Consistent Naming**: Generate Azure resource names following organizational standards with minimal configuration
- **State Persistence**: Resource names are stored in Terraform state, protecting against unintended resource recreation when naming logic changes
- **Flexible Templates**: Support for both global resources and child resources with configurable separators, prefixes, and suffixes
- **Random Suffixes**: Optional random suffixes for globally unique resource names (e.g., storage accounts)
- **Instance Numbering**: Built-in support for numbered instances with configurable padding
- **Region Functions**: Provider functions to convert between Azure region names (full, short, and CLI formats)
- **Clean Output**: Automatic removal of special characters to ensure Azure naming compliance

## Why This Provider?

When working with Azure, maintaining consistent resource naming is critical. Azure treats resource names as immutable identifiers, so changes to name generation logic can trigger destructive rebuilds that cascade across your infrastructure. By persisting names in Terraform state, this provider protects you from unintended resource recreation and the associated downtime and data loss.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the Provider

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/BHoggs/azname). Add it to your Terraform configuration:

```terraform
terraform {
  required_providers {
    azname = {
      source  = "BHoggs/azname"
      version = "~> 1.0"
    }
  }
}

provider "azname" {
  separator       = "-"
  random_length   = 4
  instance_length = 3
  clean_output    = true
}
```

### Generate Resource Names with Data Source

```terraform
data "azname_name" "resource_group" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "rg"
  location      = "westus2"
}

output "rg_name" {
  value = data.azname_name.resource_group.result
  # Example output: rg-myapp-prod-westus2-1234
}
```

### Generate Resource Names with Resource

Resources persist names in state, providing stability across runs:

```terraform
resource "azname_name" "storage" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "st"
  location      = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                = azname_name.storage.result
  resource_group_name = azurerm_resource_group.example.name
  location           = azurerm_resource_group.example.location
  # ... other configuration
}
```

### Use Region Helper Functions

```terraform
# Convert region names between formats
locals {
  region_short = provider::azname::region_short_name("West US 2")    # "wus2"
  region_cli   = provider::azname::region_cli_name("West US 2")      # "westus2"
  region_full  = provider::azname::region_full_name("westus2")       # "West US 2"
}
```

For complete documentation and examples, see the [provider documentation](https://registry.terraform.io/providers/BHoggs/azname/latest/docs).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

### Local Development

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory:

```shell
go install
```

For local development and testing, you can use the provider locally by building it and configuring Terraform to use your local build. See the [Terraform documentation on local provider development](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) for details.

```json
  dev_overrides {
    "BHoggs/azname" = "<your GOPATH>/bin"
  }
```

### Documentation

To generate or update documentation, run:

```shell
make gendocs
```

### Testing

Run the full suite of acceptance tests:

```shell
make testacc
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
