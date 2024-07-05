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

data "azname_name" "example" {}
