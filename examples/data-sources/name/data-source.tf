terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/brett/azname"
    }
  }
}

provider "azname" {
  random_length = 4
  separator     = "-"
  clean_output  = true
}

# Basic resource name generation
data "azname_name" "resource_group" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "rg"
  location      = "westus2"
}

# Output: rg-myapp-prod-westus2-1234

# Using a custom name
data "azname_name" "legacy_storage" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "st"
  custom_name   = "legacystorage123"
}

# Output: legacystorage123

# Child resource using parent name
data "azname_name" "container" {
  name          = "images"
  environment   = "prod"
  resource_type = "container"
  parent_name   = data.azname_name.legacy_storage.result
  service       = "blob"
}

# Output: legacystorage123-container-images-1234

data "scaffolding_example" "example" {
  configurable_attribute = "some-value"
}
