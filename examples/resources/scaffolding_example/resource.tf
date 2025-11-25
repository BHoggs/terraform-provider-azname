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
}

# Resource that will generate a new name when triggers change
resource "azname_name" "web_app" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "webapp"
  location      = "westus2"
  service       = "web"
  
  # Name will be regenerated if any of these change
  triggers = {
    git_sha = "abc123"
    version = "1.0.0"
  }
}

# Child resource that depends on parent
resource "azname_name" "deployment_slot" {
  name          = "blue"
  environment   = "prod"
  resource_type = "slot"
  parent_name   = azname_name.web_app.result
  
  triggers = {
    parent_name = azname_name.web_app.result
  }
}

# Using instance numbers for multiple similar resources
resource "azname_name" "vm" {
  name          = "web"
  environment   = "prod"
  resource_type = "vm"
  location      = "westus2"
  instance      = 1
}

# Using custom name for imported legacy resource
resource "azname_name" "legacy_db" {
  name          = "myapp"
  environment   = "prod"
  resource_type = "sql"
  custom_name   = "legacy-sql-db-01"
}
