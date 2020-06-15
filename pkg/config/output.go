package config

import (
	"sort"
	"strings"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/feloy/kubernetes-api-reference/pkg/outputs"
	"github.com/go-openapi/spec"
)

// OutputDocument outputs contents using output
func (o *TOC) OutputDocument(output outputs.Output) error {
	err := output.Prepare()
	if err != nil {
		return err
	}

	for p, tocPart := range o.Parts {
		err = o.OutputPart(p, tocPart, output)
		if err != nil {
			return err
		}
	}
	return nil
}

// OutputPart outputs a Part
func (o *TOC) OutputPart(i int, part *Part, output outputs.Output) error {
	outputPart, err := output.AddPart(i, part.Name)
	if err != nil {
		return err
	}

	for c, tocChapter := range part.Chapters {
		err = o.OutputChapter(c, tocChapter, outputPart)
		if err != nil {
			return err
		}
	}
	return nil
}

// OutputChapter outputs a chapter of the part
func (o *TOC) OutputChapter(i int, chapter *Chapter, outputPart outputs.Part) error {
	outputChapter, err := outputPart.AddChapter(i, chapter.Name)
	if err != nil {
		return err
	}

	if chapter.Group != nil && chapter.Version != nil {
		err = outputChapter.SetAPIVersion(GetGV(*chapter.Group, *chapter.Version))
		if err != nil {
			return err
		}
		err = outputChapter.SetGoImport(chapter.Key.GoImportPrefix())
		if err != nil {
			return err
		}
	}

	for s, tocSection := range chapter.Sections {
		err = o.OutputSection(s, tocSection, outputChapter)
		if err != nil {
			return err
		}
	}
	return nil
}

// OutputSection outputs a section of the chapter
func (o *TOC) OutputSection(i int, section *Section, outputChapter outputs.Chapter) error {
	outputSection, err := outputChapter.AddSection(i, section.Name)
	if err != nil {
		return err
	}
	err = outputSection.AddContent(section.Definition.Description)
	if err != nil {
		return err
	}

	return o.OutputProperties(section.Definition, outputSection, []string{})
}

// OutputProperties outputs the properties of a definition
func (o *TOC) OutputProperties(definition spec.Schema, outputSection outputs.Section, prefix []string) error {
	requiredProperties := definition.Required

	ordered := orderedPropertyKeys(definition.Properties)

	for _, name := range ordered {
		details := definition.Properties[name]
		property := kubernetes.NewProperty(name, details, requiredProperties)
		var linkend []string
		if property.TypeKey != nil {
			linkend = o.LinkEnds[property.TypeKey.String()]
		}
		completeName := prefix
		completeName = append(completeName, name)
		err := outputSection.AddProperty(strings.Join(completeName, "."), property, linkend, len(prefix) > 0)
		if err != nil {
			return err
		}
		if property.TypeKey != nil && len(linkend) == 0 {
			if target, found := (*o.Definitions)[property.TypeKey.String()]; found {
				o.OutputProperties(target, outputSection, append(prefix, name))
			}
		}
	}
	return nil
}

// orderedPropertyKeys returns the keys of m alphabetically ordered
func orderedPropertyKeys(m map[string]spec.Schema) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
