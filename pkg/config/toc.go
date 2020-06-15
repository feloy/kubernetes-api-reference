package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/feloy/kubernetes-api-reference/pkg/outputs/hugo"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

// TOC is the table of contents of the documentation
type TOC struct {
	Parts       []*Part `yaml:"parts"`
	LinkEnds    map[string][]string
	Definitions *spec.Definitions
}

// Part contains chapters
type Part struct {
	Name     string     `yaml:"name"`
	Chapters []*Chapter `yaml:"chapters"`
}

// Chapter contains a definition of a main resource and its associated resources and definitions
type Chapter struct {
	Name     string                 `yaml:"name"`
	Group    *kubernetes.APIGroup   `yaml:"group"`
	Version  *kubernetes.APIVersion `yaml:"version"`
	Key      kubernetes.Key         `yaml:"key"`
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
	o.LinkEnds = make(map[string][]string)

	for _, part := range o.Parts {
		for _, chapter := range part.Chapters {
			if len(chapter.Key.String()) > 0 {
				main := spec.GetDefinition(chapter.Key)
				if main != nil {
					newSection := NewSection(chapter.Name, main)
					chapter.Sections = []*Section{
						newSection,
					}
					continue
					//					o.LinkEnds[key.String()+"."+newSection.Name] = []string{part.Name, chapter.Name, newSection.Name}
				}
				return fmt.Errorf("Resource %s/%s/%s not found in spec", chapter.Group, chapter.Version.String(), kubernetes.APIKind(chapter.Name))
			}

			key, main := spec.GetResource(*chapter.Group, *chapter.Version, kubernetes.APIKind(chapter.Name), true)
			if main != nil {
				chapter.Key = key
				newSection := NewSection(chapter.Name, main)
				chapter.Sections = []*Section{
					newSection,
				}
				o.LinkEnds[key.String()+"."+newSection.Name] = []string{part.Name, chapter.Name, newSection.Name}
			} else {
				return fmt.Errorf("Resource %s/%s/%s not found in spec", chapter.Group, chapter.Version.String(), kubernetes.APIKind(chapter.Name))
			}

			suffixes := []string{"Spec", "Status", "List"}
			for _, suffix := range suffixes {
				resourceName := chapter.Name + suffix
				key, resource := spec.GetResource(*chapter.Group, *chapter.Version, kubernetes.APIKind(resourceName), true)
				if resource != nil {
					newSection := NewSection(resourceName, resource)
					chapter.Sections = append(chapter.Sections, newSection)
					o.LinkEnds[key.String()+"."+newSection.Name] = []string{part.Name, chapter.Name, newSection.Name}
				}
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
					Group:   &v.Group,
					Version: &v.Version,
					Key:     v.Key.RemoveResourceName(),
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

// ToHugo outputs documentation in Markdown format for Hugo in dir directory
func (o *TOC) ToHugo(dir string) error {
	// Test that dir is empty
	fileinfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Unable to open directory %s", dir)
	}
	if len(fileinfos) > 0 {
		return fmt.Errorf("Directory %s must be empty", dir)
	}

	hugo := hugo.NewHugo(dir)

	o.OutputDocument(hugo)
	return nil
}
