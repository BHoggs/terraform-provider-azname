terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/BHoggs/azname"
    }
  }
}

# Convert Azure region to short abbreviation format
# Useful for creating compact resource names
output "short_name_from_cli" {
  value = provider::azname::region_short_name("westus2")
  # Returns: "wus2"
}

output "short_name_from_display" {
  value = provider::azname::region_short_name("West US 2")
  # Returns: "wus2"
}

# Common use case: creating abbreviated region identifiers
# Useful when you need compact names or have length constraints
locals {
  regions = ["westus2", "eastus", "northeurope"]
  region_codes = [
    for region in local.regions :
    provider::azname::region_short_name(region)
  ]
  # Returns: ["wus2", "eus", "neu"]
}
