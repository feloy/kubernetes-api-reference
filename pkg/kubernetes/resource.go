package kubernetes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
)

// Key of the resource in the OpenAPI definition (e.g. io.k8s.api.core.v1.Pod)
type Key string

// GoImportPrefix returns the path to use for this group in go import
func (o Key) GoImportPrefix() string {
	parts := strings.Split(string(o), ".")
	return parts[1] + "." + parts[0] + "/" + strings.Join(parts[2:], "/")
}

// RemoveResourceName removes the last part of the key corresponding to the resource namz
func (o Key) RemoveResourceName() Key {
	parts := strings.Split(string(o), ".")
	return Key(strings.Join(parts[:len(parts)-1], "."))
}

// Resource represent a Kubernetes API resource
type Resource struct {
	Key        Key
	Group      APIGroup
	Version    APIVersion
	Kind       APIKind
	Definition spec.Schema

	// Replaced indicates if this version is replaced by another one
	ReplacedBy *Key
	// Documented indicates if this resource was included in the TOC
	Documented bool
}

// LessThan returns true if 'o' is a newer version than 'p'
func (o *Resource) LessThan(p *Resource) bool {
	return o.Group.Replaces(p.Group) || p.Version.LessThan(&o.Version)
}

// Replaces returns true if 'o' replaces 'p'
func (o *Resource) Replaces(p *Resource) bool {
	return o.Group.Replaces(p.Group) || o.Version.Replaces(&p.Version)
}

// Equals returns true if a resource is referenced by group/version/kind
func (o *Resource) Equals(group APIGroup, version APIVersion, kind APIKind) bool {
	return o.Group == group && o.Version.Equals(&version) && o.Kind == kind
}

// GetGV returns the group/version of a resource (used for apiVersion:)
func (o *Resource) GetGV() string {
	if o.Group == "" {
		return o.Version.String()
	}
	return fmt.Sprintf("%s/%s", o.Group, o.Version.String())
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
		for _, otherResource := range list {
			if resource.Replaces(otherResource) {
				otherResource.ReplacedBy = &resource.Key
				break
			} else if otherResource.Replaces(resource) {
				resource.ReplacedBy = &otherResource.Key
				break
			}
		}
		list = append(list, resource)
	} else {
		list = []*Resource{resource}
	}
	sort.Sort(list)
	(*o)[resource.Kind] = list
}
