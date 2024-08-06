// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNameDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr("data.azname_name.storage", "result", "unitsttesttstae852"),
				),
			},
		},
	})
}
