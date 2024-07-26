package provider

import (
	"context"
	"fmt"
	"math/rand/v2"

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

	d.config = req.ProviderData.(*aznameProviderModel)
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

	config := *d.config

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	var rng *rand.Rand
	if !state.RandomSeed.IsNull() {
		seed := uint64(state.RandomSeed.ValueInt64())
		rng = rand.New(rand.NewPCG(seed, seed))
	} else {
		rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}

	// Create and format the random suffix
	random_suffix := rng.IntN(int(config.RandomLength.ValueInt64()))
	random_suffix_str := fmt.Sprintf("%0*d", config.RandomLength.ValueInt64(), random_suffix)

	println(random_suffix_str)

	// Append global and local prefixes and suffixes
	prefixes := config.Prefixes.Elements()
	if !state.Prefixes.IsNull() {
		prefixes = append(prefixes, state.Prefixes.Elements()...)
	}

	suffixes := config.Suffixes.Elements()
	if !state.Suffixes.IsNull() {
		suffixes = append(suffixes, state.Suffixes.Elements()...)
	}

}
