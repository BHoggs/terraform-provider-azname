package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNameDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
					data "azname_name" "rg" {
						name          = "test"
						environment   = "tst"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
						prefixes      = ["unit"]
						#custom_name   = "mycustomname"
						random_seed   = 123
					}
					data "azname_name" "storage" {
						name          = "test"
						environment   = "tst"
						resource_type = "azurerm_storage_account"
						location      = "Australia East"
						prefixes      = ["unit"]
						#custom_name   = "mycustomname"
						random_seed   = 123
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azname_name.rg", "result", "unit-rg-test-tst-ae"),
					resource.TestCheckResourceAttr("data.azname_name.storage", "result", "unitsttesttstae851"),
				),
			},
		},
	})
}

func TestNameDataSourceWithProviderEnvironment(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test provider-level environment configuration
			{
				Config: `
					provider "azname" {
						random_length = 3
						environment   = "prod"
					}
					data "azname_name" "rg" {
						name          = "myapp"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azname_name.rg", "result", "rg-myapp-prod-ae"),
				),
			},
			// Test resource-level override of provider environment
			{
				Config: `
					provider "azname" {
						random_length = 3
						environment   = "prod"
					}
					data "azname_name" "rg" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azname_name.rg", "result", "rg-myapp-dev-ae"),
				),
			},
			// Test provider-level environment with global scope resource (includes random suffix)
			{
				Config: `
					provider "azname" {
						random_length = 3
						environment   = "prod"
					}
					data "azname_name" "storage" {
						name          = "myapp"
						resource_type = "azurerm_storage_account"
						location      = "Australia East"
						random_seed   = 456
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azname_name.storage", "result", "stmyappprodae382"),
				),
			},
			// Test with no environment at either level - should omit environment from name
			{
				Config: `
					provider "azname" {
						random_length = 3
					}
					data "azname_name" "rg" {
						name          = "myapp"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azname_name.rg", "result", "rg-myapp-ae"),
				),
			},
		},
	})
}
