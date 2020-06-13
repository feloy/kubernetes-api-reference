package kubernetes

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// GVKExtension represents the OpenAPI extension x-kubernetes-group-version-kind
type GVKExtension struct {
	Group   APIGroup
	Version APIVersion
	Kind    APIKind
}

// getGVKExtension returns the GVK Kubernetes extension of a definition, if found
func getGVKExtension(definition spec.Schema) (*GVKExtension, bool, error) {
	extensions := definition.Extensions
	extension, found := extensions["x-kubernetes-group-version-kind"]
	if !found {
		return nil, false, nil
	}
	gvks, ok := extension.([]interface{})
	if !ok {
		return nil, false, fmt.Errorf("x-kubernetes-group-version-kind is not an array")
	}

	if len(gvks) == 0 {
		return nil, false, nil
	}

	if len(gvks) > 1 {
		// TODO
		//t.Errorf("%s: Count of x-kubernetes-group-version-kind should be 1 but is %d", key, len(gvks))
		return nil, false, nil
	}

	gvkMap, ok := (gvks[0]).(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK")
	}

	group, ok := gvkMap["group"].(string)
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK apigroup")
	}

	version, ok := gvkMap["version"].(string)
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK apiversion")
	}

	apiversion, err := NewAPIVersion(version)
	if err != nil {
		return nil, false, fmt.Errorf("Error creating APIVersion")
	}

	kind, ok := gvkMap["kind"].(string)
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK apikind")
	}
	return &GVKExtension{
		Group:   APIGroup(group),
		Version: *apiversion,
		Kind:    APIKind(kind),
	}, true, nil
}
