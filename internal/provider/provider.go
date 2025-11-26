package provider

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &AznameProvider{}
	_ provider.ProviderWithFunctions = &AznameProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AznameProvider{
			version: version,
		}
	}
}

// AznameProvider is the provider implementation.
type AznameProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AznameProviderModel maps provider schema data to a Go type.
type AznameProviderModel struct {
	Template       types.String `tfsdk:"template"`
	TemplateChild  types.String `tfsdk:"template_child"`
	Separator      types.String `tfsdk:"separator"`
	Prefixes       types.List   `tfsdk:"prefixes"`
	Suffixes       types.List   `tfsdk:"suffixes"`
	CleanOutput    types.Bool   `tfsdk:"clean_output"`
	TrimOutput     types.Bool   `tfsdk:"trim_output"`
	RandomLength   types.Int64  `tfsdk:"random_length"`
	InstanceLength types.Int64  `tfsdk:"instance_length"`
	Environment    types.String `tfsdk:"environment"`
	Location       types.String `tfsdk:"location"`
}

// Metadata returns the provider type name.
func (p *AznameProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azname"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *AznameProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for generating standardized Azure resource names.",

		Attributes: map[string]schema.Attribute{
			"template": schema.StringAttribute{
				Optional:            true,
				Description:         "Global template for resource name generation. Default: {prefix}~{resource_type}~{workload}~{environment}~{service}~{location}{instance}{rand}~{suffix}",
				MarkdownDescription: "Global template for resource name generation. Uses ~ as a placeholder for the separator character. Can be set via `AZNAME_TEMPLATE` environment variable.",
			},
			"template_child": schema.StringAttribute{
				Optional:            true,
				Description:         "Template for child resource name generation. Default: {parent_name}~{resource_type}{instance}~{rand}",
				MarkdownDescription: "Template for child resource name generation. Uses ~ as a placeholder for the separator character. Can be set via `AZNAME_TEMPLATE_CHILD` environment variable.",
			},
			"separator": schema.StringAttribute{
				Optional:            true,
				Description:         "Character to use as separator in resource names. Default: -",
				MarkdownDescription: "Character to use as separator in resource names. Must be a single character. Can be set via `AZNAME_SEPARATOR` environment variable.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1),
				},
			},
			"prefixes": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "List of prefixes to prepend to resource names.",
				MarkdownDescription: "List of prefixes to prepend to resource names. These will be joined using the separator character. Can be set via `AZNAME_PREFIX` environment variable (comma-separated).",
			},
			"suffixes": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "List of suffixes to append to resource names.",
				MarkdownDescription: "List of suffixes to append to resource names. These will be joined using the separator character. Can be set via `AZNAME_SUFFIX` environment variable (comma-separated).",
			},
			"clean_output": schema.BoolAttribute{
				Optional:            true,
				Description:         "Remove special characters from generated names. Default: true",
				MarkdownDescription: "Remove special characters from generated names to ensure compatibility with Azure naming rules. Can be set via `AZNAME_CLEAN_OUTPUT` environment variable (1 for true, 0 for false).",
			},
			"trim_output": schema.BoolAttribute{
				Optional:            true,
				Description:         "Trim generated names to fit Azure resource length limits. Default: true",
				MarkdownDescription: "Trim generated names to fit Azure resource length limits while preserving important parts. Can be set via `AZNAME_TRIM_OUTPUT` environment variable (1 for true, 0 for false).",
			},
			"random_length": schema.Int64Attribute{
				Optional:            true,
				Description:         "Length of random suffix to append. Default: 3",
				MarkdownDescription: "Length of random suffix to append to generated names. Must be between 1 and 6. Can be set via `AZNAME_RANDOM_LENGTH` environment variable.",
				Validators: []validator.Int64{
					int64validator.Between(1, 6),
				},
			},
			"instance_length": schema.Int64Attribute{
				Optional:            true,
				Description:         "Length of instance number padding. Default: 3",
				MarkdownDescription: "Length of instance number padding in generated names. Must be between 1 and 6. Can be set via `AZNAME_INSTANCE_LENGTH` environment variable.",
				Validators: []validator.Int64{
					int64validator.Between(1, 6),
				},
			},
			"environment": schema.StringAttribute{
				Optional:            true,
				Description:         "Default environment name for all resources. Default: empty",
				MarkdownDescription: "Default environment name (e.g., dev, test, prod) to use in resource names. Can be overridden at resource/data source level. Can be set via `AZNAME_ENVIRONMENT` environment variable.",
			},
			"location": schema.StringAttribute{
				Optional:            true,
				Description:         "Default location for all resources. Default: empty",
				MarkdownDescription: "Default location (e.g., eastus, westeurope) to use in resource names. Can be overridden at resource/data source level. Can be set via `AZNAME_LOCATION` environment variable.",
			},
		},
	}
}

func (p *AznameProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config AznameProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, ok := os.LookupEnv("AZNAME_TEMPLATE")
	if !ok {
		template = "{prefix}~{resource_type}~{workload}~{environment}~{service}~{location}{instance}{rand}~{suffix}"
	}
	template_child, ok := os.LookupEnv("AZNAME_TEMPLATE_CHILD")
	if !ok {
		template_child = "{parent_name}~{resource_type}{instance}~{rand}"
	}
	separator, ok := os.LookupEnv("AZNAME_SEPARATOR")
	if !ok {
		separator = "-"
	}
	prefixes := os.Getenv("AZNAME_PREFIX")
	suffixes := os.Getenv("AZNAME_SUFFIX")
	clean_output, ok := os.LookupEnv("AZNAME_CLEAN_OUTPUT")
	if !ok {
		clean_output = "1"
	}
	trim_output, ok := os.LookupEnv("AZNAME_TRIM_OUTPUT")
	if !ok {
		trim_output = "1"
	}
	random_length, ok := os.LookupEnv("AZNAME_RANDOM_LENGTH")
	if !ok {
		random_length = "3"
	}
	instance_length, ok := os.LookupEnv("AZNAME_INSTANCE_LENGTH")
	if !ok {
		instance_length = "3"
	}
	environment, ok := os.LookupEnv("AZNAME_ENVIRONMENT")
	if !ok {
		environment = ""
	}
	location, ok := os.LookupEnv("AZNAME_LOCATION")
	if !ok {
		location = ""
	}

	// Check for required attributes, and set defaults.
	if config.Template.IsNull() {
		config.Template = types.StringValue(template)
	}
	if config.TemplateChild.IsNull() {
		config.TemplateChild = types.StringValue(template_child)
	}
	if config.Separator.IsNull() {
		config.Separator = types.StringValue(separator)
	}
	if config.Prefixes.IsNull() {
		prefixList := strings.Split(prefixes, ",")
		var attrPrefixes []attr.Value
		for _, prefix := range prefixList {
			attrPrefixes = append(attrPrefixes, types.StringValue(prefix))
		}
		config.Prefixes, diags = types.ListValue(types.StringType, attrPrefixes)
		resp.Diagnostics.Append(diags...)
	}
	if config.Suffixes.IsNull() {
		suffixList := strings.Split(suffixes, ",")
		var attrSuffixes []attr.Value
		for _, suffix := range suffixList {
			attrSuffixes = append(attrSuffixes, types.StringValue(suffix))
		}
		config.Suffixes, diags = types.ListValue(types.StringType, attrSuffixes)
		resp.Diagnostics.Append(diags...)
	}
	if config.CleanOutput.IsNull() {
		config.CleanOutput = types.BoolValue(clean_output == "1")
	}
	if config.TrimOutput.IsNull() {
		config.TrimOutput = types.BoolValue(trim_output == "1")
	}
	if config.RandomLength.IsNull() {
		randomLength, err := strconv.ParseInt(random_length, 10, 64)
		if err != nil || randomLength < 1 || randomLength > 6 {
			resp.Diagnostics.AddError("Invalid value for AZNAME_RANDOM_LENGTH", "The value must be a number between 1 and 6")
		}
		config.RandomLength = types.Int64Value(randomLength)
	}
	if config.InstanceLength.IsNull() {
		instanceLength, err := strconv.ParseInt(instance_length, 10, 64)
		if err != nil || instanceLength < 1 || instanceLength > 6 {
			resp.Diagnostics.AddError("Invalid value for AZNAME_INSTANCE_LENGTH", "The value must be a number between 1 and 6")
		}
		config.InstanceLength = types.Int64Value(instanceLength)
	}
	if config.Environment.IsNull() {
		config.Environment = types.StringValue(environment)
	}
	if config.Location.IsNull() {
		config.Location = types.StringValue(location)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = &config
	resp.DataSourceData = &config
}

// DataSources defines the data sources implemented in the provider.
func (p *AznameProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAzNameDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *AznameProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAznameResource,
	}
}

func (p *AznameProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewCliNameFunction,
		NewFullNameFunction,
		NewShortNameFunction,
	}
}
