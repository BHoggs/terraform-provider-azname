terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/hogginz/azname"
    }
  }
}

provider "azname" {
  random_length = 2
}

data "azname_name" "example" {
  name          = "test"
  environment   = "tst"
  resource_type = "azurerm_resource_group"
  location      = "Australia East"
  custom_name   = "mycustomname"
}

data "azname_name" "storage" {
  name          = "test"
  environment   = "tst"
  resource_type = "azurerm_storage_account"
  location      = "New Zealand North"
}

resource "azname_name" "storage" {
  name          = "test"
  environment   = "tst"
  resource_type = "azurerm_storage_account"
  location      = "New Zealand North"
  custom_name   = "mycustomname"

  triggers = {
    foo = "baz"
  }
}

output "name" {
  value = data.azname_name.example.result
}

output "storage" {
  value = data.azname_name.storage.result
}

output "storage_resource" {
  value = azname_name.storage.result
}

output "test_region" {
  value = provider::azname::region_cli_name("Australia East")
}