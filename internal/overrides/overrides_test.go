package overrides

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOverrides(t *testing.T) {
	// Create a temporary YAML file for testing
	tmpDir := t.TempDir()
	validYAML := `resource_slug_overrides:
  azurerm_resource_group: "resourcegroup"
  azurerm_storage_account: "storage"
region_shortname_overrides:
  eastus: "use"
  westus2: "usw2"
new_resources:
  azurerm_custom_resource:
    slug: "custom"
    min_length: 1
    max_length: 63
    scope: "resourceGroup"
    dashes: true
    lowercase: true
new_regions:
  customregion:
    cli_name: "customregion"
    full_name: "Custom Region"
    short_name: "cust"
`
	validFile := filepath.Join(tmpDir, "valid.yaml")
	if err := os.WriteFile(validFile, []byte(validYAML), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("Load valid YAML file", func(t *testing.T) {
		ovr, err := LoadOverrides(validFile)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if ovr == nil {
			t.Fatal("Expected overrides, got nil")
		}
		if len(ovr.ResourceSlugOverrides) != 2 {
			t.Errorf("Expected 2 resource slug overrides, got %d", len(ovr.ResourceSlugOverrides))
		}
		if ovr.ResourceSlugOverrides["azurerm_resource_group"] != "resourcegroup" {
			t.Errorf("Expected 'resourcegroup', got '%s'", ovr.ResourceSlugOverrides["azurerm_resource_group"])
		}
		if len(ovr.RegionShortnameOverrides) != 2 {
			t.Errorf("Expected 2 region shortname overrides, got %d", len(ovr.RegionShortnameOverrides))
		}
		if ovr.RegionShortnameOverrides["eastus"] != "use" {
			t.Errorf("Expected 'use', got '%s'", ovr.RegionShortnameOverrides["eastus"])
		}
		if len(ovr.NewResources) != 1 {
			t.Errorf("Expected 1 new resource, got %d", len(ovr.NewResources))
		}
		if ovr.NewResources["azurerm_custom_resource"].Slug != "custom" {
			t.Errorf("Expected slug 'custom', got '%s'", ovr.NewResources["azurerm_custom_resource"].Slug)
		}
		if len(ovr.NewRegions) != 1 {
			t.Errorf("Expected 1 new region, got %d", len(ovr.NewRegions))
		}
		if ovr.NewRegions["customregion"].ShortName != "cust" {
			t.Errorf("Expected shortname 'cust', got '%s'", ovr.NewRegions["customregion"].ShortName)
		}
	})

	t.Run("Load non-existent file", func(t *testing.T) {
		_, err := LoadOverrides(filepath.Join(tmpDir, "nonexistent.yaml"))
		if err == nil {
			t.Fatal("Expected error for non-existent file, got nil")
		}
	})

	t.Run("Load invalid YAML", func(t *testing.T) {
		invalidYAML := `invalid: yaml: content:`
		invalidFile := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(invalidFile, []byte(invalidYAML), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err := LoadOverrides(invalidFile)
		if err == nil {
			t.Fatal("Expected error for invalid YAML, got nil")
		}
	})
}

func TestValidateOverrides(t *testing.T) {
	t.Run("Valid new resource", func(t *testing.T) {
		ovr := &Overrides{
			NewResources: map[string]NewResourceDefinition{
				"test_resource": {
					Slug:      "test",
					MinLength: 1,
					MaxLength: 10,
					Scope:     "resourceGroup",
					Dashes:    true,
					Lowercase: false,
				},
			},
		}
		err := validateOverrides(ovr)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("New resource missing slug", func(t *testing.T) {
		ovr := &Overrides{
			NewResources: map[string]NewResourceDefinition{
				"test_resource": {
					MinLength: 1,
					MaxLength: 10,
					Scope:     "resourceGroup",
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for missing slug, got nil")
		}
	})

	t.Run("New resource missing scope", func(t *testing.T) {
		ovr := &Overrides{
			NewResources: map[string]NewResourceDefinition{
				"test_resource": {
					Slug:      "test",
					MinLength: 1,
					MaxLength: 10,
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for missing scope, got nil")
		}
	})

	t.Run("New resource invalid scope", func(t *testing.T) {
		ovr := &Overrides{
			NewResources: map[string]NewResourceDefinition{
				"test_resource": {
					Slug:      "test",
					MinLength: 1,
					MaxLength: 10,
					Scope:     "invalid",
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for invalid scope, got nil")
		}
	})

	t.Run("New resource min_length > max_length", func(t *testing.T) {
		ovr := &Overrides{
			NewResources: map[string]NewResourceDefinition{
				"test_resource": {
					Slug:      "test",
					MinLength: 20,
					MaxLength: 10,
					Scope:     "resourceGroup",
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for min_length > max_length, got nil")
		}
	})

	t.Run("Valid new region", func(t *testing.T) {
		ovr := &Overrides{
			NewRegions: map[string]NewRegionDefinition{
				"test_region": {
					CliName:   "testregion",
					FullName:  "Test Region",
					ShortName: "test",
				},
			},
		}
		err := validateOverrides(ovr)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("New region missing cli_name", func(t *testing.T) {
		ovr := &Overrides{
			NewRegions: map[string]NewRegionDefinition{
				"test_region": {
					FullName:  "Test Region",
					ShortName: "test",
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for missing cli_name, got nil")
		}
	})

	t.Run("New region missing full_name", func(t *testing.T) {
		ovr := &Overrides{
			NewRegions: map[string]NewRegionDefinition{
				"test_region": {
					CliName:   "testregion",
					ShortName: "test",
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for missing full_name, got nil")
		}
	})

	t.Run("New region missing short_name", func(t *testing.T) {
		ovr := &Overrides{
			NewRegions: map[string]NewRegionDefinition{
				"test_region": {
					CliName:  "testregion",
					FullName: "Test Region",
				},
			},
		}
		err := validateOverrides(ovr)
		if err == nil {
			t.Error("Expected error for missing short_name, got nil")
		}
	})
}

func TestDiscoverAndLoadOverrides(t *testing.T) {
	t.Run("No override file found", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		ovr, err := DiscoverAndLoadOverrides()
		if err != nil {
			t.Errorf("Expected no error when file not found, got: %v", err)
		}
		if ovr != nil {
			t.Error("Expected nil overrides when no file found, got non-nil")
		}
	})

	t.Run("Override file in current directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		validYAML := `version: "1.0"
resource_slug_overrides:
  azurerm_resource_group: "rg2"
`
		overrideFile := filepath.Join(tmpDir, "azname_overrides.yaml")
		if err := os.WriteFile(overrideFile, []byte(validYAML), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		ovr, err := DiscoverAndLoadOverrides()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if ovr == nil {
			t.Fatal("Expected overrides, got nil")
		}
		if ovr.ResourceSlugOverrides["azurerm_resource_group"] != "rg2" {
			t.Errorf("Expected 'rg2', got '%s'", ovr.ResourceSlugOverrides["azurerm_resource_group"])
		}
	})
}
