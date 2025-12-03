package resources

// Running go generate will generate the models_generated.go file
//go:generate go run gen.go

import (
	"context"
	"fmt"
	"sync"

	"terraform-provider-azname/internal/overrides"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Looks up a resource definition by its type name.
func GetResourceDefinition(resourceType string) (ResourceStructure, error) {
	resource, ok := ResourceDefinitions[resourceType]
	if !ok {
		return ResourceStructure{}, fmt.Errorf("unknown resource type: %s", resourceType)
	}
	return resource, nil
}

var overridesOnce sync.Once

// ApplyOverrides merges override configuration into the ResourceDefinitions map.
// This function is thread-safe and will only execute once using sync.Once.
func ApplyOverrides(ctx context.Context, ovr *overrides.Overrides) {
	if ovr == nil {
		return
	}

	overridesOnce.Do(func() {
		// Apply slug overrides to existing resources
		if ovr.ResourceSlugOverrides != nil {
			for resourceType, newSlug := range ovr.ResourceSlugOverrides {
				if resource, ok := ResourceDefinitions[resourceType]; ok {
					tflog.Debug(ctx, "Applying resource slug override", map[string]interface{}{
						"resource_type": resourceType,
						"old_slug":      resource.CafPrefix,
						"new_slug":      newSlug,
					})
					resource.CafPrefix = newSlug
					ResourceDefinitions[resourceType] = resource
				} else {
					tflog.Debug(ctx, "Skipping slug override for unknown resource type", map[string]interface{}{
						"resource_type": resourceType,
					})
				}
			}
		}

		// Add new resources from overrides
		if ovr.NewResources != nil {
			for resourceType, newResource := range ovr.NewResources {
				tflog.Debug(ctx, "Adding new resource type", map[string]interface{}{
					"resource_type": resourceType,
					"slug":          newResource.Slug,
					"scope":         newResource.Scope,
				})
				ResourceDefinitions[resourceType] = ResourceStructure{
					ResourceTypeName: resourceType,
					CafPrefix:        newResource.Slug,
					MinLength:        newResource.MinLength,
					MaxLength:        newResource.MaxLength,
					LowerCase:        newResource.Lowercase,
					RegEx:            "", // No cleanup regex for custom resources
					ValidationRegExp: "", // No validation regex for custom resources
					Dashes:           newResource.Dashes,
					Scope:            newResource.Scope,
				}
			}
		}
	})
}

type ResourceStructure struct {
	// Resource type name
	ResourceTypeName string `json:"name"`
	// Resource prefix as defined in the Azure Cloud Adoption Framework
	CafPrefix string `json:"slug,omitempty"`
	// MaxLength attribute define the maximum length of the name
	MinLength int `json:"min_length"`
	// MaxLength attribute define the maximum length of the name
	MaxLength int `json:"max_length"`
	// enforce lowercase
	LowerCase bool `json:"lowercase,omitempty"`
	// Regular expression to apply to the resource type
	RegEx string `json:"regex,omitempty"`
	// the Regular expression to validate the generated string
	ValidationRegExp string `json:"validatation_regex,omitempty"`
	// can the resource include dashes
	Dashes bool `json:"dashes"`
	// The scope of this name where it needs to be unique
	Scope string `json:"scope,omitempty"`
}
