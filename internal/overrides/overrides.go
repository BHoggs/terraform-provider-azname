package overrides

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Overrides represents the complete override configuration from azname_overrides.yaml
type Overrides struct {
	// Override slugs for existing resources
	ResourceSlugOverrides map[string]string `yaml:"resource_slug_overrides"`

	// Override shortnames for existing regions
	RegionShortnameOverrides map[string]string `yaml:"region_shortname_overrides"`

	// Define completely new resources not in the provider
	NewResources map[string]NewResourceDefinition `yaml:"new_resources"`

	// Define completely new regions not in the provider
	NewRegions map[string]NewRegionDefinition `yaml:"new_regions"`
}

// NewResourceDefinition defines a custom resource type with simplified schema
// (no regex validation required)
type NewResourceDefinition struct {
	// Resource prefix/slug (e.g., "rg", "st")
	Slug string `yaml:"slug"`

	// Minimum length of the generated name
	MinLength int `yaml:"min_length"`

	// Maximum length of the generated name
	MaxLength int `yaml:"max_length"`

	// Scope where the name must be unique: "global", "resourceGroup", or "parent"
	Scope string `yaml:"scope"`

	// Whether dashes are allowed in the name
	Dashes bool `yaml:"dashes"`

	// Whether the name should be lowercase
	Lowercase bool `yaml:"lowercase"`
}

// NewRegionDefinition defines a custom region with display name and short name
type NewRegionDefinition struct {
	// CLI name (e.g., "westus2")
	CliName string `yaml:"cli_name"`

	// Full display name (e.g., "West US 2")
	FullName string `yaml:"full_name"`

	// Short name used in name generation (e.g., "wus2")
	ShortName string `yaml:"short_name"`
}

// LoadOverrides loads override configuration from the specified file path
func LoadOverrides(filePath string) (*Overrides, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("overrides file not found: %s", filePath)
	}

	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read overrides file: %w", err)
	}

	// Parse YAML
	var overrides Overrides
	if err := yaml.Unmarshal(data, &overrides); err != nil {
		return nil, fmt.Errorf("failed to parse overrides file: %w", err)
	}

	// Validate the loaded overrides
	if err := validateOverrides(&overrides); err != nil {
		return nil, fmt.Errorf("invalid overrides configuration: %w", err)
	}

	return &overrides, nil
}

// DiscoverAndLoadOverrides attempts to auto-discover and load an overrides file
// Searches for ./azname_overrides.yaml in the current working directory
// Returns nil without error if no file is found (graceful degradation)
func DiscoverAndLoadOverrides() (*Overrides, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check for overrides file in current directory
	filePath := filepath.Join(cwd, "azname_overrides.yaml")
	if _, err := os.Stat(filePath); err == nil {
		// File found, try to load it
		return LoadOverrides(filePath)
	}

	// No file found - this is not an error, just return nil
	return nil, nil
}

// validateOverrides performs basic validation on the override configuration
func validateOverrides(o *Overrides) error {
	// Validate new resource definitions
	for name, resource := range o.NewResources {
		// Check required fields
		if resource.Slug == "" {
			return fmt.Errorf("new_resources[%s]: slug is required", name)
		}
		if resource.Scope == "" {
			return fmt.Errorf("new_resources[%s]: scope is required (must be 'global', 'resourceGroup', or 'parent')", name)
		}

		// Validate scope value
		if resource.Scope != "global" && resource.Scope != "resourceGroup" && resource.Scope != "parent" {
			return fmt.Errorf("new_resources[%s]: scope must be 'global', 'resourceGroup', or 'parent', got '%s'", name, resource.Scope)
		}

		// Validate length constraints
		if resource.MinLength < 0 {
			return fmt.Errorf("new_resources[%s]: min_length cannot be negative", name)
		}
		if resource.MaxLength < 1 {
			return fmt.Errorf("new_resources[%s]: max_length must be at least 1", name)
		}
		if resource.MinLength > resource.MaxLength {
			return fmt.Errorf("new_resources[%s]: min_length (%d) cannot be greater than max_length (%d)", name, resource.MinLength, resource.MaxLength)
		}
	}

	// Validate new region definitions
	for name, region := range o.NewRegions {
		// Check required fields
		if region.CliName == "" {
			return fmt.Errorf("new_regions[%s]: cli_name is required", name)
		}
		if region.FullName == "" {
			return fmt.Errorf("new_regions[%s]: full_name is required", name)
		}
		if region.ShortName == "" {
			return fmt.Errorf("new_regions[%s]: short_name is required", name)
		}

		// Warn if shortname is very long (but don't fail)
		if len(region.ShortName) > 10 {
			fmt.Fprintf(os.Stderr, "Warning: new_regions[%s]: short_name '%s' is longer than 10 characters, which may cause naming issues\n", name, region.ShortName)
		}
	}

	return nil
}
