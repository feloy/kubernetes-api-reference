package kubernetes

import (
	"fmt"
	"strings"

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

	GVToKey GVToKeyMap
}

// GVToKeyMap maps Kubernetes resource Group/Version with Spec Definition key (without Kind)
// e.g. GVToKey["v1"]: "io.k8s.api.core.v1"
type GVToKeyMap map[string]string

// Add adds a new match between key and resource GV
func (o GVToKeyMap) Add(key string, resource *Resource) {
	parts := strings.Split(key, ".")
	if len(parts) == 0 {
		return
	}
	subkey := strings.Join(parts[0:len(parts)-1], ".")
	gv := resource.GetGV()
	if _, found := o[gv]; !found {
		o[gv] = subkey
	}
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
// and maps definitions keys to Resources GVs
func (o *Spec) getResources() error {
	o.Resources = &ResourceMap{}
	o.GVToKey = GVToKeyMap{}

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

		resource := &Resource{
			Key:        key,
			Group:      APIGroup(group),
			Version:    *apiversion,
			Kind:       APIKind(kind),
			Definition: definition,
		}
		o.Resources.Add(resource)
		o.GVToKey.Add(key, resource)
	}
	return nil
}

// GetResource returns the resource referenced by group/version/kind, or nil if not found
func (o *Spec) GetResource(group APIGroup, version APIVersion, kind APIKind, markAsDocumented bool) *spec.Schema {
	// Search on K8s resources
	for k, resources := range *o.Resources {
		if k == kind {
			for r, resource := range resources {
				if resource.Equals(group, version, kind) {
					if markAsDocumented {
						(*o.Resources)[k][r].Documented = true
					}
					return &resource.Definition
				}
			}
		}
	}

	// Get on definitions
	gvRes := Resource{
		Group:   group,
		Version: version,
	}
	gvk := o.GVToKey[gvRes.GetGV()] + "." + kind.String()
	if def, found := o.Swagger.Definitions[gvk]; found {
		return &def
	}

	return nil
}
