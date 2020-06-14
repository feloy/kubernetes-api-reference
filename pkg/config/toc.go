package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/feloy/kubernetes-api-reference/pkg/formats/markdown"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/feloy/kubernetes-api-reference/pkg/outputs/hugo"
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
	Key      kubernetes.Key
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
			key, main := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(chapter.Name), true)
			if main != nil {
				chapter.Key = key
				chapter.Sections = []*Section{
					NewSection(chapter.Name, main),
				}
			} else {
				return fmt.Errorf("Resource %s/%s/%s not found in spec", chapter.Group, chapter.Version.String(), kubernetes.APIKind(chapter.Name))
			}

			suffixes := []string{"Spec", "Status", "List"}
			for _, suffix := range suffixes {
				resourceName := chapter.Name + suffix
				_, resource := spec.GetResource(chapter.Group, chapter.Version, kubernetes.APIKind(resourceName), true)
				if resource != nil {
					chapter.Sections = append(chapter.Sections, NewSection(resourceName, resource))
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
					Group:   v.Group,
					Version: v.Version,
					Key:     v.Key,
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
	err = hugo.AddIndex("", map[string]interface{}{
		"title": "Resources",
	})
	if err != nil {
		return fmt.Errorf("Error writing index file in %s: %s", dir, err)
	}

	for p, part := range o.Parts {
		partname, err := hugo.AddPart(part.Name)
		if err != nil {
			return fmt.Errorf("Error creating part %s: %s", part.Name, err)
		}

		err = hugo.AddIndex(partname, map[string]interface{}{
			"title":       part.Name,
			"draft":       false,
			"collapsible": true,
			"weight":      p + 1,
		})

		for c, chapter := range part.Chapters {
			chaptername, err := hugo.AddChapter(partname, chapter.Name, map[string]interface{}{
				"title":       chapter.Name,
				"draft":       false,
				"collapsible": false,
				"weight":      c + 1,
			})
			if err != nil {
				return fmt.Errorf("Error creating chapter %s/%s: %s", part.Name, chapter.Name, err)
			}

			err = hugo.AddContent(partname, chaptername, markdown.Code("apiVersion: "+GetGV(chapter.Group, chapter.Version)))
			if err != nil {
				return fmt.Errorf("Error adding GV for chapter %s/%s: %s", part.Name, chapter.Name, err)
			}

			err = hugo.AddContent(partname, chaptername, markdown.Code("import \""+chapter.Key.GoImportPrefix()+"\""))
			if err != nil {
				return fmt.Errorf("Error adding Go Import for chapter %s/%s: %s", part.Name, chapter.Name, err)
			}

			for _, section := range chapter.Sections {
				hugo.AddSection(partname, chaptername, section.Name)
				hugo.AddContent(partname, chaptername, section.Definition.Description)
			}
		}
	}

	return nil
}
