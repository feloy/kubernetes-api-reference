package config

import (
	"fmt"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/go-openapi/spec"
)

// Chapter contains a definition of a main resource and its associated resources and definitions
type Chapter struct {
	Name     string                 `yaml:"name"`
	Group    *kubernetes.APIGroup   `yaml:"group"`
	Version  *kubernetes.APIVersion `yaml:"version"`
	Key      kubernetes.Key         `yaml:"key"`
	Sections []*Section
}

func (o *Chapter) isResource() bool {
	return o.Group != nil && o.Version != nil
}

func (o *Chapter) populate(part *Part, toc *TOC, thespec *kubernetes.Spec) error {
	var main *spec.Schema

	if o.isResource() {
		if o.Key, main = thespec.GetResource(*o.Group, *o.Version, kubernetes.APIKind(o.Name), true); main == nil {
			return fmt.Errorf("Resource %s/%s/%s not found in spec", o.Group, o.Version.String(), kubernetes.APIKind(o.Name))
		}
	} else {
		if main = thespec.GetDefinition(o.Key); main == nil {
			return fmt.Errorf("Resource %s/%s/%s not found in spec", o.Group, o.Version.String(), kubernetes.APIKind(o.Name))
		}
	}

	newSection := NewSection(o.Name, main)
	o.Sections = []*Section{
		newSection,
	}
	toc.LinkEnds.Add(o.Key, []string{part.Name, o.Name + "-" + o.Version.String(), newSection.Name})
	toc.DocumentedDefinitions[o.Key] = []string{o.Name}

	if o.isResource() {
		suffixes := []string{"Spec", "Status"}
		for _, suffix := range suffixes {
			resourceName := o.Name + suffix
			gvRes := kubernetes.Resource{
				GVKExtension: kubernetes.GVKExtension{
					Group:   *o.Group,
					Version: *o.Version,
				},
			}
			keys := thespec.GVToKey[gvRes.GetGV()]
			for _, key := range keys {
				resourceKey := kubernetes.Key(key + "." + resourceName)
				o.addDefinition(resourceName, resourceKey, part, toc, thespec)
			}
		}

		suffixes = []string{"List"}
		for _, suffix := range suffixes {
			resourceName := o.Name + suffix
			key, resource := thespec.GetResource(*o.Group, *o.Version, kubernetes.APIKind(resourceName), true)
			if resource != nil {
				newSection := NewSection(resourceName, resource)
				o.Sections = append(o.Sections, newSection)
				toc.LinkEnds.Add(key, []string{part.Name, o.Name + "-" + o.Version.String(), newSection.Name})
				toc.DocumentedDefinitions[key] = []string{resourceName}
			}
		}

	} else {

		suffixes := []string{"Spec", "Status", "List"}
		for _, suffix := range suffixes {
			resourceName := o.Name + suffix
			resourceKey := kubernetes.Key(o.Key.String() + suffix)
			o.addDefinition(resourceName, resourceKey, part, toc, thespec)
		}
	}
	return nil
}

func (o *Chapter) addDefinition(resourceName string, resourceKey kubernetes.Key, part *Part, toc *TOC, thespec *kubernetes.Spec) {
	resource := thespec.GetDefinition(resourceKey)
	if resource != nil {
		newSection := NewSection(resourceName, resource)
		o.Sections = append(o.Sections, newSection)
		toc.LinkEnds.Add(resourceKey, []string{part.Name, o.Name + "-" + o.Version.String(), newSection.Name})
		toc.DocumentedDefinitions[resourceKey] = []string{resourceName}
	}
}
