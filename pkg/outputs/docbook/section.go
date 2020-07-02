package docbook

import (
	"strings"

	"github.com/feloy/kubernetes-api-reference/pkg/formats/dbxml"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	x "github.com/shabbyrobe/xmlwriter"
)

type DocbookSection struct {
	w *x.Writer
}

// AddContent adds content to the output
func (o DocbookSection) AddContent(s string) error {
	for _, part := range strings.Split(s, "\n") {
		o.w.StartElem(dbxml.ElemWithText("para", part))
		o.w.EndElem("para")
	}
	return nil
}

// AddTypeDefinition adds the definition of a type to the output
func (o DocbookSection) AddTypeDefinition(s string) error {
	for _, part := range strings.Split(s, "\n") {
		o.w.StartElem(x.Elem{Name: "para"})
		o.w.StartElem(dbxml.ElemWithText("emphasis", part))
		o.w.EndElem("emphasis")
		o.w.EndElem("para")
	}
	return nil
}

// StartPropertyList starts the list of properties
func (o DocbookSection) StartPropertyList() error {
	return o.w.StartElem(x.Elem{Name: "variablelist"})
}

// AddProperty adds a property to the list of properties
func (o DocbookSection) AddProperty(name string, property *kubernetes.Property, linkend []string, indent bool) error {
	o.w.StartElem(x.Elem{Name: "varlistentry"})
	o.w.StartElem(x.Elem{Name: "term"})
	o.w.StartElem(dbxml.ElemWithText("varname", name))
	o.w.EndElem("varname")
	o.w.WriteText(" (")
	o.w.StartElem(x.Elem{Name: "emphasis"})
	if len(linkend) > 0 {
		o.w.StartElem(x.Elem{Name: "link"})
		o.w.WriteAttr(x.Attr{Name: "linkend", Value: escapeName(linkend[1] + "." + linkend[2])})
	}
	o.w.WriteText(property.Type)
	if len(linkend) > 0 {
		o.w.EndElem("link")
	}
	o.w.EndElem("emphasis")
	o.w.WriteText(")")
	if property.Required {
		o.w.WriteText(", required")
	}
	o.w.EndElem("term")
	o.w.StartElem(x.Elem{Name: "listitem"})

	var patches []x.Writable
	if property.MergeStrategyKey != nil && property.RetainKeysStrategy {
		patches = []x.Writable{x.Text("Patch strategies: retainKeys, merge on key "), dbxml.ElemWithText("varname", *property.MergeStrategyKey)}
	} else if property.MergeStrategyKey != nil {
		patches = []x.Writable{x.Text("Patch strategy: merge on key "), dbxml.ElemWithText("varname", *property.MergeStrategyKey)}
	} else if property.RetainKeysStrategy {
		patches = []x.Writable{x.Text("Patch strategy: retainKeys")}
	}

	if len(patches) > 0 {
		o.w.StartElem(x.Elem{Name: "para"})
		o.w.StartElem(dbxml.ElemWithContent("emphasis", patches))
		o.w.EndElem("emphasis")
		o.w.EndElem("para")
	}

	parts := strings.Split(property.Description, "\n")
	for _, part := range parts {
		o.w.StartElem(dbxml.ElemWithText("para", part))
		o.w.EndElem("para")
	}
	return nil
}

// EndProperty ends a property
func (o DocbookSection) EndProperty() error {
	o.w.EndElem("listitem")
	o.w.EndElem("varlistentry")
	return nil
}

// EndPropertyList ends the list of properties
func (o DocbookSection) EndPropertyList() error {
	return o.w.EndElem("variablelist")
}
