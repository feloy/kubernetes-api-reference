package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

// TOC is the table of contents of the documentation
type TOC struct {
	Parts []*Part `yaml:"parts"`
}

// Part contains chapters
type Part struct {
	Name     string     `yaml:"name"`
	Chapters []*Chapter `yaml:"chapters"`
}

// Chapter contains a definition of a main resource and its associated resources and definitions
type Chapter struct {
	Name     string                `yaml:"name"`
	Group    kubernetes.APIGroup   `yaml:"group"`
	Version  kubernetes.APIVersion `yaml:"version"`
	Sections []*Section
}

// Section contains a definition of a Kind for a given Group/Version
type Section struct {
	Name       string
	Definition spec.Schema
}

// NewSection returns a Section
func NewSection(name string, definition *spec.Schema) *Section {
	return &Section{
		Name:       name,
		Definition: *definition,
	}
}

// LoadTOC loads a config file containing the TOC definition
func LoadTOC(filename string) (*TOC, error) {
	var result TOC

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// PopulateAssociates adds sections to the chapters found in the spec
func (o *TOC) PopulateAssociates(spec *kubernetes.Spec) error {
	for _, part := range o.Parts {
		for _, chapter := range part.Chapters {
			main := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name), true)
			if main != nil {
				chapter.Sections = []*Section{
					NewSection(chapter.Name, main),
				}
			} else {
				return fmt.Errorf("Resource %s/%s/%s not found in spec", chapter.Group, chapter.Version.String(), kubernetes.APIKind(chapter.Name))
			}

			specResource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name+"Spec"), true)
			if specResource != nil {
				chapter.Sections = append(chapter.Sections, NewSection(chapter.Name+"Spec", specResource))
			}

			statusResource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name+"Status"), true)
			if statusResource != nil {
				chapter.Sections = append(chapter.Sections, NewSection(chapter.Name+"Status", statusResource))
			}

			listResource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name+"List"), true)
			if listResource != nil {
				chapter.Sections = append(chapter.Sections, NewSection(chapter.Name+"List", listResource))
			}
		}
	}
	return nil
}

// AddOtherResources adds not documented and not replaced resources to a new Part
func (o *TOC) AddOtherResources(spec *kubernetes.Spec) {
	part := &Part{}
	part.Name = "Other Resources"
	part.Chapters = []*Chapter{}

	for _, resource := range *spec.Resources {
		for _, v := range resource {
			if v.ReplacedBy == nil && !v.Documented {
				part.Chapters = append(part.Chapters, &Chapter{
					Name:    v.Kind.String(),
					Group:   v.Group,
					Version: v.Version,
				})
			}
		}
	}
	sort.Slice(part.Chapters, func(i, j int) bool {
		return part.Chapters[i].Name < part.Chapters[j].Name
	})
	if len(part.Chapters) > 0 {
		o.Parts = append(o.Parts, part)
	}
}

// ToMarkdown writes a Markdown representation of the TOC
func (o *TOC) ToMarkdown(w io.Writer) {
	for _, part := range o.Parts {
		fmt.Fprintf(w, "\n## %s\n", part.Name)
		for _, chapter := range part.Chapters {
			fmt.Fprintf(w, "### %s\n", chapter.Name)
			for _, section := range chapter.Sections {
				fmt.Fprintf(w, "#### %s\n", section.Name)
			}
		}
	}
}

// GetGV returns the group/version for a resource and version (used for apiVersion:)
func GetGV(group kubernetes.APIGroup, version kubernetes.APIVersion) string {
	if group == "" {
		return version.String()
	}
	return fmt.Sprintf("%s/%s", group, version.String())
}
