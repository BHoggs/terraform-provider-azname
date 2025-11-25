terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/brett/azname"
    }
  }
}

# Provider-level configuration sets defaults that apply to all resources
# These settings (separator, random_length, prefixes, suffixes) can be
# configured here to avoid repeating them in each resource
provider "azname" {
  random_length = 4
  separator     = "-"
}

# Basic resource name generation
# Required: name, environment, resource_type
# Resource-specific attributes like location and service should be set here
resource "azname_name" "web_app" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "webapp"
  location      = "westus2" # Resource-specific, should be set at resource level
  service       = "web"     # Resource-specific, should be set at resource level

  # Name will be regenerated if any of these triggers change
  # Useful for forcing new names when deployment parameters change
  triggers = {
    git_sha = "abc123"
    version = "1.0.0"
  }
}

# Child resource that depends on parent
# Uses parent_name to create hierarchical naming
resource "azname_name" "deployment_slot" {
  name          = "blue"
  environment   = "prod"
  resource_type = "slot"
  parent_name   = azname_name.web_app.result

  # Trigger regeneration when parent name changes
  triggers = {
    parent_name = azname_name.web_app.result
  }
}

# Using instance numbers for multiple similar resources
# The instance number will be formatted according to provider's instance_length setting
resource "azname_name" "vm" {
  name          = "web"
  environment   = "prod"
  resource_type = "vm"
  location      = "westus2"
  instance      = 1
}

# Using custom name for imported legacy resource
# When custom_name is set, it overrides the generated name entirely
resource "azname_name" "legacy_db" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "sql"
  custom_name   = "legacy-sql-db-01"
}

# Example showing resource-level override of provider settings
# You can override separator or other provider settings per-resource if needed
resource "azname_name" "special_storage" {
  name          = "data"
  environment   = "prod"
  resource_type = "st"
  location      = "eastus"
  separator     = "" # Override provider separator for this specific resource

  # Resource-level prefixes/suffixes are also supported
  prefixes = ["corp"]
  suffixes = ["primary"]
}
