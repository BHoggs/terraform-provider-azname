package provider

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"regexp"
	"strings"

	"terraform-provider-azname/internal/regions"
	"terraform-provider-azname/internal/resources"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Helper function to convert a Terraform list to a Go slice.
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

func generateName(ctx context.Context, state aznameDataSourceModel, config aznameProviderModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	resourceType, ok := resources.ResourceDefinitions[state.ResourceType.ValueString()]
	if !ok {
		diags.AddAttributeError(path.Root("resource_type"), "unknown resource type", fmt.Sprintf("Unknown resource type: %s", state.ResourceType.ValueString()))
		return "", diags
	}

	var randomSuffixString string
	if resourceType.Scope == "global" {
		var rng *rand.Rand
		if !state.RandomSeed.IsNull() {
			seed := uint64(state.RandomSeed.ValueInt64())
			rng = rand.New(rand.NewPCG(seed, seed))
		} else {
			rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
		}
		randomLength := int(config.RandomLength.ValueInt64())
		randomSuffix := rng.IntN(int(math.Pow10(randomLength)))
		randomSuffixString = fmt.Sprintf("%0*d", randomLength, randomSuffix)

	}

	prefixes, err := convertFromTfList[string](ctx, config.Prefixes)
	if err != nil {
		diags.AddError("Error extracting prefixes", err.Error())
		return "", diags
	}

	if !state.Prefixes.IsNull() {
		prefixes, err = convertFromTfList[string](ctx, state.Prefixes)
		if err != nil {
			diags.AddError("Error extracting prefixes", err.Error())
			return "", diags
		}
	}

	suffixes, err := convertFromTfList[string](ctx, config.Suffixes)
	if err != nil {
		diags.AddError("Error extracting suffixes", err.Error())
		return "", diags
	}

	if !state.Suffixes.IsNull() {
		suffixes, err = convertFromTfList[string](ctx, state.Suffixes)
		if err != nil {
			diags.AddError("Error extracting suffixes", err.Error())
			return "", diags
		}
	}

	regionShortName, err := regions.GetRegionByAnyName(state.Location.ValueString())
	if err != nil {
		diags.AddAttributeError(path.Root("location"), "unknown region", fmt.Sprintf("Unknown region: %s", state.Location.ValueString()))
		return "", diags
	}

	var instanceString string
	if !state.Instance.IsNull() {
		instanceString = fmt.Sprintf("%0*d", config.InstanceLength.ValueInt64(), state.Instance.ValueInt64())
	}

	replacer := strings.NewReplacer(
		"{prefix}", strings.Join(prefixes, "~"),
		"{parent_name}", state.ParentName.ValueString(),
		"{resource_type}", resourceType.CafPrefix,
		"{workload}", state.Name.ValueString(),
		"{service}", state.Service.ValueString(),
		"{environment}", state.Environment.ValueString(),
		"{location}", regionShortName.ShortName,
		"{suffix}", strings.Join(suffixes, "~"),
		"{instance}", instanceString,
		"{rand}", randomSuffixString,
	)

	result := config.Template.ValueString()

	if !state.ParentName.IsNull() {
		result = config.TemplateChild.ValueString()
	}

	result = replacer.Replace(result)

	result = regexp.MustCompile(`~{2,}`).ReplaceAllString(result, "~")
	result = strings.Trim(result, "~")

	separator := config.Separator.ValueString()
	if !state.Separator.IsNull() {
		separator = state.Separator.ValueString()
	}

	if !resourceType.Dashes {
		separator = ""
	}

	result = strings.ReplaceAll(result, "~", separator)

	// clean output
	if config.CleanOutput.ValueBool() {
		result = regexp.MustCompile(resourceType.RegEx).ReplaceAllString(result, "")
	}

	// trim output to length
	if config.TrimOutput.ValueBool() {
		// runes are more reliable than bytes for trimming
		runes := []rune(result)
		trimLength := min(len(runes), resourceType.MaxLength)
		result = string(runes[:trimLength])
	}

	// validate the output
	r := regexp.MustCompile(resourceType.ValidationRegExp)
	if !r.MatchString(result) {
		diags.AddError("Generated name failed validation", fmt.Sprintf("Generated name %q failed validation against %q", result, resourceType.ValidationRegExp))
	}

	return result, diags
}
