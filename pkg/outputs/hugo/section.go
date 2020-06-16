package hugo

import (
	"fmt"

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
