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

func NewFullNameFunction() function.Function {
	return FullNameFunction{}
}

func NewShortNameFunction() function.Function {
	return ShortNameFunction{}
}

type CliNameFunction struct{}
type FullNameFunction struct{}
type ShortNameFunction struct{}

func (r CliNameFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "region_cli_name"
}

func (r FullNameFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "region_full_name"
}

func (r ShortNameFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "region_short_name"
}

func (r CliNameFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Azure Region CLI Name",
		MarkdownDescription: `Gets the Azure CLI name for a region. 
This function takes a region name in any format (full name, short name, or CLI name) and returns the CLI name.`,
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "region",
				MarkdownDescription: "Region full name (e.g. 'West US 2'), short name (e.g. 'wus2'), or CLI name (e.g. 'westus2')",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r FullNameFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Azure Region Full Name",
		MarkdownDescription: `Gets the Azure full name for a region.
This function takes a region name in any format (full name, short name, or CLI name) and returns the full display name.`,
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "region",
				MarkdownDescription: "Region full name (e.g. 'West US 2'), short name (e.g. 'wus2'), or CLI name (e.g. 'westus2')",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r ShortNameFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Azure Region Short Name",
		MarkdownDescription: `Gets a CAF recommended short name for a region.
This function takes a region name in any format (full name, short name, or CLI name) and returns the Cloud Adoption Framework (CAF) recommended short name.`,
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "region",
				MarkdownDescription: "Region full name (e.g. 'West US 2'), short name (e.g. 'wus2'), or CLI name (e.g. 'westus2')",
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

func (r FullNameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
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

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, region.FullName))
}

func (r ShortNameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
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

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, region.ShortName))
}
