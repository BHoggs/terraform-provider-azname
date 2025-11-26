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
	Result       types.String `tfsdk:"result"`
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
		Description: "Resource for generating standardized Azure resource names following naming conventions.",

		Attributes: map[string]schema.Attribute{
			"result": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "The generated resource name following the configured template pattern.",
				MarkdownDescription: "The generated resource name following the configured template pattern.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The workload or application name to use in the resource name.",
				MarkdownDescription: "The workload or application name to use in the resource name.",
			},
			"environment": schema.StringAttribute{
				Required:            true,
				Description:         "The environment name (e.g., dev, test, prod) to use in the resource name.",
				MarkdownDescription: "The environment name (e.g., dev, test, prod) to use in the resource name.",
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
			"triggers": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
				Description:         "Map of values that should trigger a new name to be generated when changed.",
				MarkdownDescription: "Map of values that should trigger a new name to be generated when changed. Common triggers include version numbers or Git commit hashes.",
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

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
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

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
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
