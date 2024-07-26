// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-azname/internal/regions"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = CliNameFunction{}
)

func NewCliNameFunction() function.Function {
	return CliNameFunction{}
}

type CliNameFunction struct{}

func (r CliNameFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "region_cli_name"
}

func (r CliNameFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Azure Region CLI Name",
		MarkdownDescription: "Gets the Azure CLI name for a region",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "region",
				MarkdownDescription: "Region full, short, or CLI name",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r CliNameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var inputRegion string

	resp.Error = req.Arguments.Get(ctx, &inputRegion)
	if resp.Error != nil {
		return
	}

	region, err := regions.GetRegionByAnyName(inputRegion)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("region not found: %s", inputRegion))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, region.CliName))
}
