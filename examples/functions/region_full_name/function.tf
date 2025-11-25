terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/brett/azname"
    }
  }
}

# Convert Azure region CLI name to full display name
# Useful for displaying user-friendly region names
output "full_name_from_cli" {
  value = provider::azname::region_full_name("westus2")
  # Returns: "West US 2"
}

output "full_name_from_short" {
  value = provider::azname::region_full_name("wus2")
  # Returns: "West US 2"
}

# Common use case: displaying region names in outputs or UI
resource "azname_name" "example" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "rg"
  location      = "westus2"
}

output "friendly_location" {
  description = "User-friendly region name for documentation"
  value       = provider::azname::region_full_name("westus2")
}
