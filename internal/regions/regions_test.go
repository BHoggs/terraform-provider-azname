package regions

import (
	"testing"
)

func TestGetRegionByAnyName(t *testing.T) {
	// Test case 1: Valid short name
	shortName := "ne"
	expectedRegion := "northeurope"
	region, err := GetRegionByAnyName(shortName)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if region.CliName != expectedRegion {
		t.Errorf("Expected region %v, got %v", expectedRegion, region.CliName)
	}

	// Test case 2: Valid CLI name
	cliName := "eastus"
	expectedRegion = "East US"
	region, err = GetRegionByAnyName(cliName)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if region.FullName != expectedRegion {
		t.Errorf("Expected region %v, got %v", expectedRegion, region.FullName)
	}

	// Test case 3: Valid full name
	fullName := "Australia East"
	expectedRegion = "ae"
	region, err = GetRegionByAnyName(fullName)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if region.ShortName != expectedRegion {
		t.Errorf("Expected region %v, got %v", expectedRegion, region.ShortName)
	}

	// Test case 4: Invalid name
	invalidName := "invalid"
	region, err = GetRegionByAnyName(invalidName)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if region != nil {
		t.Errorf("Expected nil region, got %v", region)
	}
}
