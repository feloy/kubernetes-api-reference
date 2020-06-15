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
func (o Section) AddProperty(name string, property *kubernetes.Property, linkend []string) error {
	required := ""
	if property.Required {
		required = ", required"
	}

	link := ""
	var title string
	if property.TypeKey != nil {
		link = " [" + property.TypeKey.String() + "]"
		if len(linkend) > 0 {
			link = o.hugo.LinkEnd(linkend, property.Type)
		}
		title = fmt.Sprintf("**%s** (%s)%s", name, link, required)
	} else {
		title = fmt.Sprintf("**%s** (%s%s)%s", name, property.Type, link, required)
	}

	return o.hugo.addListEntry(o.part.name, o.chapter.name, title, property.Description)
}
