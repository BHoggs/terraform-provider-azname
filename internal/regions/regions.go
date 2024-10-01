package regions

import (
	"errors"
	"strings"
)

type region struct {
	CliName      string
	FullName     string
	ShortName    string
	PairedRegion *string
}

func stringPtr(s string) *string {
	return &s
}

// Get-AzLocation | select Location, DisplayName, @{l="Pair"; e={$_.PairedRegion[0].Name}}
// Short names: https://github.com/Azure/terraform-azurerm-caf-enterprise-scale/blob/main/modules/connectivity/locals.geo_codes.tf.json
var regionsList = []region{
	{"asia", "Asia", "asia", nil},
	{"asiapacific", "Asia Pacific", "apac", nil},
	{"australia", "Australia", "aus", nil},
	{"australiacentral", "Australia Central", "acl", stringPtr("australiacentral2")},
	{"australiacentral2", "Australia Central 2", "acl2", stringPtr("australiacentral")},
	{"australiaeast", "Australia East", "ae", stringPtr("australiasoutheast")},
	{"australiasoutheast", "Australia Southeast", "ase", stringPtr("australiaeast")},
	{"brazil", "Brazil", "bra", nil},
	{"brazilsouth", "Brazil South", "brs", stringPtr("southcentralus")},
	{"brazilsoutheast", "Brazil Southeast", "bse", stringPtr("brazilsouth")},
	{"canada", "Canada", "can", nil},
	{"canadacentral", "Canada Central", "cnc", stringPtr("canadaeast")},
	{"canadaeast", "Canada East", "cne", stringPtr("canadacentral")},
	{"centralindia", "Central India", "inc", stringPtr("southindia")},
	{"centralus", "Central US", "cus", stringPtr("eastus2")},
	{"centraluseuap", "Central US EUAP", "ccy", stringPtr("eastus2euap")},
	{"eastasia", "East Asia", "ea", stringPtr("southeastasia")},
	{"eastus", "East US", "eus", stringPtr("westus")},
	{"eastus2", "East US 2", "eus2", stringPtr("centralus")},
	{"eastus2euap", "East US 2 EUAP", "ecy", stringPtr("centraluseuap")},
	{"europe", "Europe", "eu", nil},
	{"france", "France", "fra", nil},
	{"francecentral", "France Central", "frc", stringPtr("francesouth")},
	{"francesouth", "France South", "frs", stringPtr("francecentral")},
	{"germany", "Germany", "ger", nil},
	{"germanynorth", "Germany North", "gn", stringPtr("germanywestcentral")},
	{"germanywestcentral", "Germany West Central", "gwc", stringPtr("germanynorth")},
	{"global", "Global", "glob", nil},
	{"india", "India", "ind", nil},
	{"israel", "Israel", "isr", nil},
	{"israelcentral", "Israel Central", "ilc", nil},
	{"italy", "Italy", "ita", nil},
	{"italynorth", "Italy North", "itn", nil},
	{"japan", "Japan", "jap", nil},
	{"japaneast", "Japan East", "jpe", stringPtr("japanwest")},
	{"japanwest", "Japan West", "jpw", stringPtr("japaneast")},
	{"korea", "Korea", "kor", nil},
	{"koreacentral", "Korea Central", "krc", stringPtr("koreasouth")},
	{"koreasouth", "Korea South", "krs", stringPtr("koreacentral")},
	{"mexicocentral", "Mexico Central", "mexc", nil},
	{"newzealandnorth", "New Zealand North", "nzn", nil},
	{"northcentralus", "North Central US", "ncus", stringPtr("southcentralus")},
	{"northeurope", "North Europe", "ne", stringPtr("westeurope")},
	{"norway", "Norway", "nor", nil},
	{"norwayeast", "Norway East", "nwe", stringPtr("norwaywest")},
	{"norwaywest", "Norway West", "nww", stringPtr("norwayeast")},
	{"poland", "Poland", "pol", nil},
	{"polandcentral", "Poland Central", "polc", nil},
	{"qatar", "Qatar", "qat", nil},
	{"qatarcentral", "Qatar Central", "qac", nil},
	{"singapore", "Singapore", "sgp", nil},
	{"southafrica", "South Africa", "saf", nil},
	{"southafricanorth", "South Africa North", "san", stringPtr("southafricawest")},
	{"southafricawest", "South Africa West", "saw", stringPtr("southafricanorth")},
	{"southcentralus", "South Central US", "scus", stringPtr("northcentralus")},
	{"southeastasia", "Southeast Asia", "sea", stringPtr("eastasia")},
	{"southindia", "South India", "ins", stringPtr("centralindia")},
	{"spaincentral", "Spain Central", "spnc", nil},
	{"sweden", "Sweden", "swe", nil},
	{"swedencentral", "Sweden Central", "sdc", stringPtr("swedensouth")},
	{"switzerland", "Switzerland", "swi", nil},
	{"switzerlandnorth", "Switzerland North", "szn", stringPtr("switzerlandwest")},
	{"switzerlandwest", "Switzerland West", "szw", stringPtr("switzerlandnorth")},
	{"uaecentral", "UAE Central", "uac", stringPtr("uaenorth")},
	{"uaenorth", "UAE North", "uan", stringPtr("uaecentral")},
	{"uksouth", "UK South", "uks", stringPtr("ukwest")},
	{"ukwest", "UK West", "ukw", stringPtr("uksouth")},
	{"unitedstates", "United States", "us", nil},
	{"westcentralus", "West Central US", "wcus", stringPtr("westus2")},
	{"westeurope", "West Europe", "we", stringPtr("northeurope")},
	{"westindia", "West India", "inw", stringPtr("southindia")},
	{"westus", "West US", "wus", stringPtr("eastus")},
	{"westus2", "West US 2", "wus2", stringPtr("westcentralus")},
	{"westus3", "West US 3", "wus3", stringPtr("eastus")},
}

// GetRegionByShortName returns a region by its short name.
func GetRegionByShortName(shortName string) (*region, error) {
	for _, r := range regionsList {
		if strings.EqualFold(r.ShortName, shortName) {
			return &r, nil
		}
	}
	return nil, errors.New("region not found")
}

// GetRegionByCliName returns a region by its CLI name.
func GetRegionByCliName(cliName string) (*region, error) {
	for _, r := range regionsList {
		if strings.EqualFold(r.CliName, cliName) {
			return &r, nil
		}
	}
	return nil, errors.New("region not found")
}

// GetRegionByFullName returns a region by its full name.
func GetRegionByFullName(fullName string) (*region, error) {
	for _, r := range regionsList {
		if strings.EqualFold(r.FullName, fullName) {
			return &r, nil
		}
	}
	return nil, errors.New("region not found")
}

func GetRegionByAnyName(name string) (*region, error) {
	r, err := GetRegionByShortName(name)
	if err == nil {
		return r, nil
	}
	r, err = GetRegionByCliName(name)
	if err == nil {
		return r, nil
	}
	return GetRegionByFullName(name)
}
