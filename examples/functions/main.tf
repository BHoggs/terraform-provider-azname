terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/brett/azname"
    }
  }
}

provider "azname" {}

# Using the region functions
output "cli_name" {
  value = provider::azname::region_cli_name("West US 2")
  # Outputs: westus2
}

output "full_name" {
  value = provider::azname::region_full_name("westus2")
  # Outputs: West US 2
}

output "short_name" {
  value = provider::azname::region_short_name("westus2")
  # Outputs: wus2
}

# Example of chaining functions
output "chained_example" {
  value = provider::azname::region_short_name(provider::azname::region_full_name("westus2"))
  # Takes westus2 -> West US 2 -> wus2
}