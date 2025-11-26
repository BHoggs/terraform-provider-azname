terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/BHoggs/azname"
    }
  }
}

# Provider configuration with global settings
# These settings (separator, random_length, clean_output, etc.) apply to all data sources
provider "azname" {
  random_length = 4
  separator     = "-"
  clean_output  = true
}

# Basic resource name generation using data source
# Data sources are useful for generating names without managing state
# Required attributes: name, environment, resource_type
data "azname_name" "resource_group" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "rg"
  location      = "westus2" # Resource-specific attribute
}

# Output: rg-myapp-prod-westus2-1234

# Using a custom name to override generation
# Useful when you need to reference existing resources with specific names
data "azname_name" "legacy_storage" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "st"
  custom_name   = "legacystorage123"
}

# Output: legacystorage123

# Child resource using parent name
# Demonstrates hierarchical naming for related resources
data "azname_name" "container" {
  name          = "images"
  environment   = "prod"
  resource_type = "container"
  parent_name   = data.azname_name.legacy_storage.result
  service       = "blob" # Resource-specific attribute
}

# Output: legacystorage123-container-images-1234
