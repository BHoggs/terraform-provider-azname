terraform {
  required_providers {
    azname = {
      source = "BHoggs/azname"
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
