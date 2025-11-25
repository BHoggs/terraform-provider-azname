terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/brett/azname"
    }
  }
}

# Convert Azure region display name to CLI format
# Useful for standardizing region names in configurations
output "cli_name_from_display" {
  value = provider::azname::region_cli_name("West US 2")
  # Returns: "westus2"
}

output "cli_name_from_short" {
  value = provider::azname::region_cli_name("wus2")
  # Returns: "westus2"
}

# Common use case: ensuring region names are in CLI format
locals {
  user_specified_region = "East US"
  normalized_region     = provider::azname::region_cli_name(local.user_specified_region)
}
