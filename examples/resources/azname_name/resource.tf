terraform {
  required_providers {
    azname = {
      source = "BHoggs/azname"
    }
  }
}

# Provider-level configuration sets defaults that apply to all resources
# These settings (separator, random_length, prefixes, suffixes, environment) can be
# configured here to avoid repeating them in each resource
provider "azname" {
  environment   = "prod"
  random_length = 4
  separator     = "-"
}

# Basic resource name generation
# Required: name, resource_type
# Resource-specific attributes like location and service should be set here
resource "azname_name" "web_app" {
  name          = "myapp"
  resource_type = "azurerm_linux_web_app"
  location      = "westus2" # Resource-specific, should be set at resource level
  service       = "web"     # Resource-specific, should be set at resource level

  # Name will be regenerated if any of these triggers change
  # Useful for forcing new names when deployment parameters change
  triggers = {
    key = "value"
  }
}

# Child resource that depends on parent
# Uses parent_name to create hierarchical naming
resource "azname_name" "deployment_slot" {
  name          = "blue"
  resource_type = "azurerm_function_app_slot"
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
  resource_type = "azurerm_virtual_machine"
  location      = "westus2"
  instance      = 1
}

# Using custom name for imported legacy resource
# When custom_name is set, it overrides the generated name entirely
resource "azname_name" "legacy_db" {
  name          = "myapp"
  resource_type = "azurerm_mssql_database"
  custom_name   = "legacy-sql-db-01"
}

# Example showing resource-level override of provider settings
# You can override environment, separator, or other provider settings per-resource if needed
resource "azname_name" "special_storage" {
  name          = "data"
  environment   = "dev" # Override provider environment for this specific resource
  resource_type = "azurerm_storage_account"
  location      = "eastus"
  separator     = "" # Override provider separator for this specific resource

  # Resource-level prefixes/suffixes are also supported
  prefixes = ["corp"]
  suffixes = ["primary"]
}

# Storage account without random_seed - random suffix shown as (known after apply)
resource "azname_name" "storage_dynamic" {
  name          = "data"
  resource_type = "azurerm_storage_account"
  location      = "eastus"
  # Result will be something like "stdata123" but shown as (known after apply) in plan
}

# Storage account with random_seed - random suffix shown in plan
resource "azname_name" "storage_deterministic" {
  name          = "data"
  resource_type = "azurerm_storage_account"
  location      = "eastus"
  random_seed   = 12345
  # Result will always be "stdata897" and shown in plan output
}
