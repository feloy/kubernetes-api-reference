package kubernetes

import (
	"sort"

	"github.com/go-openapi/spec"
)

// Resource represent a Kubernetes API resource
type Resource struct {
	// The key of the resource in the OpenAPI definition (e.g. io.k8s.api.core.v1.Pod)
	Key        string
	Group      APIGroup
	Version    APIVersion
	Kind       APIKind
	Definition *spec.Schema
}

// LessThan returns true if 'o' is an older version than 'p'
func (o *Resource) LessThan(p *Resource) bool {
	return o.Group.Replaces(p.Group) || p.Version.LessThan(&o.Version)
}

// ResourceList is the list of resources for a given Kind
type ResourceList []*Resource

func (a ResourceList) Len() int           { return len(a) }
func (a ResourceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ResourceList) Less(i, j int) bool { return a[i].LessThan(a[j]) }

// ResourceMap contains a map of resources, classified by Kind
type ResourceMap map[APIKind]ResourceList

// Add a resource to the resource list
func (o *ResourceMap) Add(resource *Resource) {
	list, ok := (*o)[resource.Kind]
	if ok {
		list = append(list, resource)
	} else {
		list = []*Resource{resource}
	}
	sort.Sort(list)
	(*o)[resource.Kind] = list
}
