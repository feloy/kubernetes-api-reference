package openapi_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/feloy/kubernetes-api-reference/pkg/openapi"
)

func Test_OpenAPISpecV118(t *testing.T) {
	filename := "../../api/v1.18.json"
	doc, _ := openapi.LoadOpenAPISpec(filename)

	resources := kubernetes.ResourceMap{}

	for k, def := range doc.Definitions {
		extensions := def.Extensions
		extension, found := extensions["x-kubernetes-group-version-kind"]
		if !found {
			continue
		}
		gvks, ok := extension.([]interface{})
		if !ok {
			t.Errorf("%s: x-kubernetes-group-version-kind is not an array", k)
		}
		if len(gvks) > 1 {
			t.Errorf("%s: Count of x-kubernetes-group-version-kind should be 1 but is %d", k, len(gvks))
			continue
		}
		gvk, ok := (gvks[0]).(map[string]interface{})
		if !ok {
			t.Errorf("%s: Error getting GVK", k)
		}
		group, ok := gvk["group"].(string)
		if !ok {
			t.Errorf("%s: Error getting GVK apigroup", k)
		}
		version, ok := gvk["version"].(string)
		if !ok {
			t.Errorf("%s: Error getting GVK apiversion", k)
		}
		apiversion, err := kubernetes.NewAPIVersion(version)
		if err != nil {
			t.Errorf("%s: Error creating APIVersion", k)
		}
		kind, ok := gvk["kind"].(string)
		if !ok {
			t.Errorf("%s: Error getting GVK apikind", k)
		}
		resources.Add(&kubernetes.Resource{
			Key:        k,
			Group:      kubernetes.APIGroup(group),
			Version:    *apiversion,
			Kind:       kubernetes.APIKind(kind),
			Definition: &def,
		})
	}

	i := 0
	keys := make([]string, len(resources))
	for k := range resources {
		keys[i] = k.String()
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		rs := resources[kubernetes.APIKind(k)]
		fmt.Printf("%s\n", k)
		for _, r := range rs {
			fmt.Printf("\t%s/%s\n", r.Group, r.Version.String())
		}
	}
}
