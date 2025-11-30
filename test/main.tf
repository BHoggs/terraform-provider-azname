terraform {
  required_providers {
    azname = {
      source  = "BHoggs/azname"
      version = "0.2.1"
    }
  }
}

locals {
  names_version = "5"
}

provider "azname" {
  #environment     = "prd"
  location        = "New Zealand North"
  random_length   = 3
  instance_length = 3
  prefixes        = ["hgg"]
}

resource "azname_name" "storage_account" {
  resource_type = "azurerm_storage_account"
  name          = "val"

  triggers = {
    names_version = local.names_version
  }
}

resource "azname_name" "storage_with_seed" {
  resource_type = "azurerm_storage_account"
  name          = "seeded"
  random_seed   = 12345
}

resource "azname_name" "resource_group" {
  resource_type = "azurerm_resource_group"
  name          = "validation"
  service       = "mgmt"
  triggers = {
    names_version = local.names_version
  }
}

output "storage_account" {
  value = azname_name.storage_account.result
}

output "storage_with_seed" {
  value = azname_name.storage_with_seed.result
}

output "resource_group" {
  value = azname_name.resource_group.result
}