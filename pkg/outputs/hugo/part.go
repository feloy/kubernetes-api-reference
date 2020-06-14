package hugo

import (
	"fmt"

	"github.com/feloy/kubernetes-api-reference/pkg/outputs"
)

// Part of a Hugo output
// implements the outputs.Part interface
type Part struct {
	hugo *Hugo
	name string
}

// AddChapter adds a chapter to the Part
func (o Part) AddChapter(i int, name string) (outputs.Chapter, error) {
	chaptername, err := o.hugo.addChapter(o.name, name, map[string]interface{}{
		"title":       name,
		"draft":       false,
		"collapsible": false,
		"weight":      i + 1,
	})
	if err != nil {
		return Chapter{}, fmt.Errorf("Error creating chapter %s/%s: %s", o.name, name, err)
	}

	return Chapter{
		hugo: o.hugo,
		part: &o,
		name: chaptername,
	}, nil
}
