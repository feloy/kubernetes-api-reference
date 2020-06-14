package kubernetes

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// Property represents a property of a definition
type Property struct {
	Name        string
	Type        string
	TypeKey     *Key
	Description string
	Required    bool
}

// NewProperty returns a new Property from its swagger definition
func NewProperty(name string, details spec.Schema, required []string) (*Property, error) {
	typ, key := getTypeNameAndKey(details)
	result := Property{
		Name:        name,
		Type:        typ,
		TypeKey:     key,
		Description: details.Description,
	}
	result.Required = isRequired(name, required)

	return &result, nil
}

// isRequired returns true if name appears in the required array
func isRequired(name string, required []string) bool {
	for _, req := range required {
		if req == name {
			return true
		}
	}
	return false
}

// getTypeName returns the display name of a Schema.
// This is the api kind for definitions and the type for
// primitive types.
func getTypeNameAndKey(s spec.Schema) (string, *Key) {
	// Get the reference for complex types
	if isDefinition(s) {
		key := Key(s.SchemaProps.Ref.GetPointer().String())
		return key.ResourceName(), &key
	}

	// Recurse if type is array
	if isArray(s) {
		typ, key := getTypeNameAndKey(*s.Items.Schema)
		return fmt.Sprintf("[]%s", typ), key
	}

	// Get the value for primitive types
	if len(s.Type) > 0 {
		return fmt.Sprintf("%s", s.Type[0]), nil
	}

	panic(fmt.Errorf("No type found for object %v", s))
}

// isDefinition returns true if Schema is a complex type that should have a Definition
func isDefinition(s spec.Schema) bool {
	return len(s.SchemaProps.Ref.GetPointer().String()) > 0
}

// isArray returns true if the type is an array type
func isArray(s spec.Schema) bool {
	return len(s.Type) > 0 && s.Type[0] == "array"
}
