package resources

// Running go generate will generate the models_generated.go file
//go:generate go run gen.go

import "fmt"

// Looks up a resource definition by its type name
func GetResourceDefinition(resourceType string) (ResourceStructure, error) {
	resource, ok := ResourceDefinitions[resourceType]
	if !ok {
		return ResourceStructure{}, fmt.Errorf("unknown resource type: %s", resourceType)
	}
	return resource, nil
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
