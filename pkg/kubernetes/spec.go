package kubernetes

import (
	"fmt"

	"github.com/feloy/kubernetes-api-reference/pkg/openapi"
	"github.com/go-openapi/spec"
)

// Spec represents the Kubernetes API Specification
type Spec struct {
	// Swagger is the openAPI representation of the k8s spec
	// populated by calling getSwagger
	Swagger *spec.Swagger

	// Resources is the list of K8s resources
	// populated by calling getResources
	Resources *ResourceMap
}

// NewSpec creates a new Spec from a K8s spec file
func NewSpec(filename string) (*Spec, error) {
	spec := &Spec{}
	err := spec.getSwagger(filename)
	if err != nil {
		return nil, err
	}
	err = spec.getResources()
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// GetSwagger populates the swagger representation of the Spec in file `filename`
func (o *Spec) getSwagger(filename string) error {
	var err error
	o.Swagger, err = openapi.LoadOpenAPISpec(filename)
	return err
}

// GetResources populates the resources defined in the spec
func (o *Spec) getResources() error {
	o.Resources = &ResourceMap{}

	for key, definition := range o.Swagger.Definitions {
		extensions := definition.Extensions
		extension, found := extensions["x-kubernetes-group-version-kind"]
		if !found {
			continue
		}
		gvks, ok := extension.([]interface{})
		if !ok {
			return fmt.Errorf("%s: x-kubernetes-group-version-kind is not an array", key)
		}

		if len(gvks) > 1 {
			// TODO
			//t.Errorf("%s: Count of x-kubernetes-group-version-kind should be 1 but is %d", key, len(gvks))
			continue
		}

		gvk, ok := (gvks[0]).(map[string]interface{})
		if !ok {
			return fmt.Errorf("%s: Error getting GVK", key)
		}

		group, ok := gvk["group"].(string)
		if !ok {
			return fmt.Errorf("%s: Error getting GVK apigroup", key)
		}

		version, ok := gvk["version"].(string)
		if !ok {
			return fmt.Errorf("%s: Error getting GVK apiversion", key)
		}

		apiversion, err := NewAPIVersion(version)
		if err != nil {
			return fmt.Errorf("%s: Error creating APIVersion", key)
		}

		kind, ok := gvk["kind"].(string)
		if !ok {
			return fmt.Errorf("%s: Error getting GVK apikind", key)
		}

		o.Resources.Add(&Resource{
			Key:        key,
			Group:      APIGroup(group),
			Version:    *apiversion,
			Kind:       APIKind(kind),
			Definition: &definition,
		})
	}
	return nil
}
