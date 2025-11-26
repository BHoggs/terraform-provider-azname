terraform {
  required_providers {
    azname = {
      source = "BHoggs/azname"
    }
  }
}

# Provider configuration with global settings
# These settings (separator, random_length, clean_output, environment, etc.) apply to all data sources
provider "azname" {
  environment   = "prod"
  random_length = 4
  separator     = "-"
  clean_output  = true
}

# Basic resource name generation using data source
# Data sources are useful for generating names without managing state
# Required attributes: name, resource_type
# Environment is inherited from provider configuration
data "azname_name" "resource_group" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  location      = "westus2" # Resource-specific attribute
}

# Output: rg-myapp-prod-westus2-1234

# Using a custom name to override generation
# Useful when you need to reference existing resources with specific names
data "azname_name" "legacy_storage" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  custom_name   = "legacystorage123"
}

# Output: legacystorage123

# Child resource using parent name
# Demonstrates hierarchical naming for related resources
data "azname_name" "container" {
  name          = "images"
  resource_type = "azurerm_storage_container"
  parent_name   = data.azname_name.legacy_storage.result
  service       = "blob" # Resource-specific attribute
}

# Output: legacystorage123-container-images-1234
