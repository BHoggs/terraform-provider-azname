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
	_ datasource.DataSource              = &aznameDataSource{}
	_ datasource.DataSourceWithConfigure = &aznameDataSource{}
)

func AzNameDataSource() datasource.DataSource {
	return &aznameDataSource{}
}

type aznameDataSource struct {
	config *aznameProviderModel
}

type aznameDataSourceModel struct {
	Name         types.String `tfsdk:"name"`
	Environment  types.String `tfsdk:"environment"`
	CustomName   types.String `tfsdk:"custom_name"`
	ResourceType types.String `tfsdk:"resource_type"`
	Prefixes     types.List   `tfsdk:"prefixes"`
	Suffixes     types.List   `tfsdk:"suffixes"`
	Separator    types.String `tfsdk:"separator"`
	RandomSeed   types.Int64  `tfsdk:"random_seed"`
	Location     types.String `tfsdk:"location"`
	Instance     types.Int64  `tfsdk:"instance"`
	Service      types.String `tfsdk:"service"`
	ParentName   types.String `tfsdk:"parent_name"`
	Result       types.String `tfsdk:"result"`
}

func (d *aznameDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*aznameProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aznameProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.config = config
}

func (d *aznameDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

func (d *aznameDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"environment": schema.StringAttribute{
				Required: true,
			},
			"custom_name": schema.StringAttribute{
				Optional: true,
			},
			"resource_type": schema.StringAttribute{
				Required: true,
			},
			"prefixes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"suffixes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"separator": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1),
				},
			},
			"random_seed": schema.Int64Attribute{
				Optional: true,
			},
			"location": schema.StringAttribute{
				Optional: true,
			},
			"instance": schema.Int64Attribute{
				Optional: true,
			},
			"service": schema.StringAttribute{
				Optional: true,
			},
			"parent_name": schema.StringAttribute{
				Optional: true,
			},
			"result": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *aznameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state aznameDataSourceModel
	var result string

	config := *d.config

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// If a custom_name is provided, use that as the result
	if !state.CustomName.IsNull() {
		result = state.CustomName.ValueString()
	} else {
		res, diags := generateName(ctx, state, config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		result = res
	}

	state.Result = types.StringValue(result)
	resp.State.Set(ctx, state)
}

func convertFromTfList[T any](ctx context.Context, list types.List) ([]T, error) {
	var result []T
	elements := list.Elements()

	for _, element := range elements {
		var value T
		tfValue, err := element.ToTerraformValue(ctx)
		if err != nil {
			return nil, err
		}
		err = tfValue.As(&value)
		if err != nil {
			return nil, err
		}

		result = append(result, value)
	}

	return result, nil
}
