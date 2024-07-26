terraform {
  required_providers {
    azname = {
      source = "registry.terraform.io/hogginz/azname"
    }
  }
}

provider "azname" {
  random_length = 3
}

data "azname_name" "example" {
  name        = "test"
  environment = "tst"
  resource_type = "test"
}

output "test_region" {
  value = provider::azname::region_cli_name("Australia East")
}