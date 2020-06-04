package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
			main := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name))
			if main != nil {
				chapter.Sections = []*Section{
					NewSection(chapter.Name, main),
				}
			} else {
				return fmt.Errorf("Resource %s/%s/%s not found in spec", chapter.Group, chapter.Version.String(), kubernetes.APIKind(chapter.Name))
			}

			specResource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name+"Spec"))
			if specResource != nil {
				chapter.Sections = append(chapter.Sections, NewSection(chapter.Name+"Spec", specResource))
			}

			statusResource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name+"Status"))
			if statusResource != nil {
				chapter.Sections = append(chapter.Sections, NewSection(chapter.Name+"Status", statusResource))
			}

			listResource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name+"List"))
			if listResource != nil {
				chapter.Sections = append(chapter.Sections, NewSection(chapter.Name+"List", listResource))
			}
		}
	}
	return nil
}

// ToMarkdown writes in w a Markdown representation of the TOC
func (o *TOC) ToMarkdown(w io.Writer) {
	for _, part := range o.Parts {
		fmt.Fprintf(w, "\n## %s\n", part.Name)
		for _, chapter := range part.Chapters {
			fmt.Fprintf(w, "### %s\n", chapter.Name)
			for _, section := range chapter.Sections {
				fmt.Fprintf(w, "#### %s\n", section.Name)
				fmt.Fprintf(w, "%s\n", section.Definition.Description)
			}
		}
	}
}
