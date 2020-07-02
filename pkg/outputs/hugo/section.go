package hugo

import (
	"fmt"

	"github.com/feloy/kubernetes-api-reference/pkg/formats/markdown"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

// Section of a Hugo output
// implements the outputs.Section interface
type Section struct {
	hugo    *Hugo
	part    *Part
	chapter *Chapter
}

// AddContent adds content to a section
func (o Section) AddContent(s string) error {
	return o.hugo.addContent(o.part.name, o.chapter.name, s)
}

// AddTypeDefinition adds the definition of a type to a section
func (o Section) AddTypeDefinition(s string) error {
	return o.hugo.addContent(o.part.name, o.chapter.name, markdown.Italic(s))
}

// StartPropertyList starts the list of properties
func (o Section) StartPropertyList() error {
	return nil
}

// AddProperty adds a property to the section
func (o Section) AddProperty(name string, property *kubernetes.Property, linkend []string, indent bool) error {
	indentLevel := 0
	if indent {
		indentLevel++
	}
	required := ""
	if property.Required {
		required = ", required"
	}

	link := ""
	var title string
	if property.TypeKey != nil {
		link = property.Type
		if len(linkend) > 0 {
			link = o.hugo.LinkEnd(linkend, property.Type)
		}
		title = fmt.Sprintf("**%s** (%s)%s", name, link, required)
	} else {
		title = fmt.Sprintf("**%s** (%s%s)%s", name, property.Type, link, required)
	}

	description := property.Description
	var patches string
	if property.MergeStrategyKey != nil && property.RetainKeysStrategy {
		patches = fmt.Sprintf("Patch strategies: retainKeys, merge on key `%s`", *property.MergeStrategyKey)
	} else if property.MergeStrategyKey != nil {
		patches = fmt.Sprintf("Patch strategy: merge on key `%s`", *property.MergeStrategyKey)
	} else if property.RetainKeysStrategy {
		patches = "Patch strategy: retainKeys"
	}

	if len(patches) > 0 {
		description = "*" + patches + "*\n" + description
	}
	return o.hugo.addListEntry(o.part.name, o.chapter.name, title, description, indentLevel)
}

// EndProperty ends a property
func (o Section) EndProperty() error {
	return nil
}

// EndPropertyList ends the list of properties
func (o Section) EndPropertyList() error {
	return nil
}
