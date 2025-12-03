package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &AznameDataSource{}
	_ datasource.DataSourceWithConfigure = &AznameDataSource{}
)

func NewAzNameDataSource() datasource.DataSource {
	return &AznameDataSource{}
}

type AznameDataSource struct {
	config *AznameProviderModel
}

type AznameDataSourceModel struct {
	AznameNameModel
}

func (d *AznameDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*AznameProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *AznameProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.config = config
}

func (d *AznameDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

func (d *AznameDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for generating standardized Azure resource names following naming conventions.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the data source, same as result.",
				MarkdownDescription: "ID of the data source, same as result.",
			},
			"result": schema.StringAttribute{
				Computed:            true,
				Description:         "The generated resource name following the configured template pattern.",
				MarkdownDescription: "The generated resource name following the configured template pattern.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The workload or application name to use in the resource name.",
				MarkdownDescription: "The workload or application name to use in the resource name.",
			},
			"environment": schema.StringAttribute{
				Optional:            true,
				Description:         "The environment name (e.g., dev, test, prod) to use in the resource name.",
				MarkdownDescription: "The environment name (e.g., dev, test, prod) to use in the resource name. Defaults to provider-level environment if not set.",
			},
			"custom_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Override the generated name with a custom value. Useful for legacy or imported resources.",
				MarkdownDescription: "Override the generated name with a custom value. Useful for legacy or imported resources.",
			},
			"resource_type": schema.StringAttribute{
				Required:            true,
				Description:         "The Azure resource type abbreviation (e.g., rg for resource group, kv for key vault).",
				MarkdownDescription: "The Azure resource type abbreviation (e.g., rg for resource group, kv for key vault).",
			},
			"prefixes": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "List of prefixes to prepend to the resource name.",
				MarkdownDescription: "List of prefixes to prepend to the resource name. These will be joined using the separator character.",
			},
			"suffixes": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "List of suffixes to append to the resource name.",
				MarkdownDescription: "List of suffixes to append to the resource name. These will be joined using the separator character.",
			},
			"separator": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1),
				},
				Description:         "Character to use as separator in the resource name. Defaults to provider's separator setting.",
				MarkdownDescription: "Character to use as separator in the resource name. Must be a single character. Defaults to provider's separator setting.",
			},
			"random_seed": schema.Int64Attribute{
				Optional:            true,
				Description:         "Seed value for random suffix generation. Use this to get consistent random values.",
				MarkdownDescription: "Seed value for random suffix generation. Use this to get consistent random values.",
			},
			"location": schema.StringAttribute{
				Optional:            true,
				Description:         "Azure region where the resource will be deployed.",
				MarkdownDescription: "Azure region where the resource will be deployed. Will be included in the name if specified in the template.",
			},
			"instance": schema.Int64Attribute{
				Optional:            true,
				Description:         "Instance number for the resource. Used when deploying multiple instances.",
				MarkdownDescription: "Instance number for the resource. Used when deploying multiple instances of the same resource type.",
			},
			"service": schema.StringAttribute{
				Optional:            true,
				Description:         "Service or component identifier within the workload.",
				MarkdownDescription: "Service or component identifier within the workload (e.g., web, api, worker).",
			},
			"parent_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Name of the parent resource for child resources.",
				MarkdownDescription: "Name of the parent resource. Required when generating names for child resources.",
			},
		},
	}
}

func (d *AznameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AznameDataSourceModel
	var result string

	config := *d.config

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If a custom_name is provided, use that as the result
	if !state.CustomName.IsNull() {
		state.Result = state.CustomName
		state.ID = state.CustomName
		resp.State.Set(ctx, state)

		return
	}

	result, diags := GenerateName(ctx, state.AznameNameModel, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Result = types.StringValue(result)
	state.ID = types.StringValue(result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
