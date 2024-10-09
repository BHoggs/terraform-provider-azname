// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AznameResource{}
var _ resource.ResourceWithImportState = &AznameResource{}

func NewAznameResource() resource.Resource {
	return &AznameResource{}
}

type AznameResource struct {
	config *AznameProviderModel
}

// These are shared between the resource and data source implementations.
type AznameNameModel struct {
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

type AznameResourceModel struct {
	AznameNameModel
	Triggers types.Map `tfsdk:"triggers"`
}

func (r *AznameResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

func (r *AznameResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"triggers": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *AznameResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

	r.config = config
}

func (r *AznameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state AznameResourceModel
	var result string

	config := *r.config

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If a custom_name is provided, use that as the result
	if !state.CustomName.IsNull() {
		state.Result = state.CustomName
		resp.State.Set(ctx, state)

		return
	}

	result, diags := GenerateName(ctx, state.AznameNameModel, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Result = types.StringValue(result)
	resp.State.Set(ctx, state)
}

func (r *AznameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This is a no-op because the resource is computed.
}

func (r *AznameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state AznameResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AznameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This is a no-op because the resource is computed.
}

func (r *AznameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("result"), req, resp)
}
