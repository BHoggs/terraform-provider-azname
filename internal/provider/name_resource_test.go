package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestNameResource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
					resource "azname_name" "rg" {
						name          = "test"
						environment   = "tst"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
						prefixes      = ["unit"]
						random_seed   = 123
					}
					resource "azname_name" "storage" {
						name          = "test"
						environment   = "tst"
						resource_type = "azurerm_storage_account"
						location      = "Australia East"
						prefixes      = ["unit"]
						random_seed   = 123
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.rg", "result", "unit-rg-test-tst-ae"),
					resource.TestCheckResourceAttr("azname_name.storage", "result", "unitsttesttstae851"),
				),
			},
			// Re-apply with different attributes to verify result doesn't change
			{
				Config: providerConfig + `
					resource "azname_name" "rg" {
						name          = "changed"
						environment   = "prod"
						resource_type = "azurerm_resource_group"
						location      = "West US"
						prefixes      = ["different"]
						random_seed   = 999
					}
					resource "azname_name" "storage" {
						name          = "changed"
						environment   = "prod"
						resource_type = "azurerm_storage_account"
						location      = "West US"
						prefixes      = ["different"]
						random_seed   = 999
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Results persist from first step despite attribute changes
					resource.TestCheckResourceAttr("azname_name.rg", "result", "unit-rg-test-tst-ae"),
					resource.TestCheckResourceAttr("azname_name.storage", "result", "unitsttesttstae851"),
				),
			},
		},
	})
}

func TestNameResource_PlanKnown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test that result is known during plan (not "known after apply")
			{
				Config: `
					provider "azname" {
						random_length = 3
					}
					resource "azname_name" "test" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
					}
					`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify the result is known in the plan (not computed after apply)
						plancheck.ExpectKnownValue(
							"azname_name.test",
							tfjsonpath.New("result"),
							knownvalue.StringExact("rg-myapp-dev-ae"),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.test", "result", "rg-myapp-dev-ae"),
				),
			},
			// Test with global scope resource (includes random suffix)
			{
				Config: `
					provider "azname" {
						random_length = 3
					}
					resource "azname_name" "storage" {
						name          = "myapp"
						environment   = "prod"
						resource_type = "azurerm_storage_account"
						location      = "East US"
						random_seed   = 999
					}
					`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify the result is known in the plan with random suffix
						plancheck.ExpectKnownValue(
							"azname_name.storage",
							tfjsonpath.New("result"),
							knownvalue.StringExact("stmyappprodeus415"),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.storage", "result", "stmyappprodeus415"),
				),
			},
		},
	})
}

func TestNameResource_CustomName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test custom_name override
			{
				Config: providerConfig + `
					resource "azname_name" "custom" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
						custom_name   = "my-custom-name"
					}
					`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify custom name is used in the plan
						plancheck.ExpectKnownValue(
							"azname_name.custom",
							tfjsonpath.New("result"),
							knownvalue.StringExact("my-custom-name"),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.custom", "result", "my-custom-name"),
					resource.TestCheckResourceAttr("azname_name.custom", "custom_name", "my-custom-name"),
				),
			},
		},
	})
}

func TestNameResource_Triggers(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with triggers
			{
				Config: providerConfig + `
					locals {
						version = "1.0"
					}
					resource "azname_name" "with_triggers" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
						triggers = {
							version = local.version
						}
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.with_triggers", "result", "azname-rg-myapp-dev-ae"),
					resource.TestCheckResourceAttr("azname_name.with_triggers", "triggers.version", "1.0"),
				),
			},
			// Update triggers - should cause replacement
			{
				Config: providerConfig + `
					locals {
						version = "2.0"
					}
					resource "azname_name" "with_triggers" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
						triggers = {
							version = local.version
						}
					}
					`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify resource will be replaced
						plancheck.ExpectResourceAction("azname_name.with_triggers", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.with_triggers", "result", "azname-rg-myapp-dev-ae"),
					resource.TestCheckResourceAttr("azname_name.with_triggers", "triggers.version", "2.0"),
				),
			},
		},
	})
}

func TestNameResource_ProviderDefaults(t *testing.T) {
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
					resource "azname_name" "rg" {
						name          = "myapp"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
					}
					`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"azname_name.rg",
							tfjsonpath.New("result"),
							knownvalue.StringExact("rg-myapp-prod-ae"),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azname_name.rg", "result", "rg-myapp-prod-ae"),
				),
			},
			// Test resource-level override of provider environment
			{
				Config: `
					provider "azname" {
						random_length = 3
						environment   = "prod"
					}
					resource "azname_name" "rg" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						location      = "Australia East"
					}
					`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"azname_name.rg",
							tfjsonpath.New("result"),
							knownvalue.StringExact("rg-myapp-prod-ae"),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					// Result persists from step 1 despite environment override
					resource.TestCheckResourceAttr("azname_name.rg", "result", "rg-myapp-prod-ae"),
				),
			},
		},
	})
}

func TestNameResource_Validation(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid separator (more than 1 character)
			{
				Config: providerConfig + `
					resource "azname_name" "invalid" {
						name          = "myapp"
						environment   = "dev"
						resource_type = "azurerm_resource_group"
						separator     = "---"
					}
					`,
				ExpectError: regexp.MustCompile("Attribute separator string length must be at most 1"),
			},
		},
	})
}
