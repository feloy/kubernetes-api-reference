package config_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/config"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/feloy/kubernetes-api-reference/pkg/outputs"
)

type FakeOutput struct{}

func (o FakeOutput) Prepare() error                                   { return nil }
func (o FakeOutput) AddPart(i int, name string) (outputs.Part, error) { return FakePart{}, nil }

type FakePart struct{}

func (o FakePart) AddChapter(i int, name string) (outputs.Chapter, error) { return FakeChapter{}, nil }

type FakeChapter struct{}

func (o FakeChapter) SetAPIVersion(s string) error { return nil }
func (o FakeChapter) SetGoImport(s string) error   { return nil }
func (o FakeChapter) AddSection(i int, name string) (outputs.Section, error) {
	return FakeSection{}, nil
}

type FakeSection struct{}

func (o FakeSection) AddContent(s string) error                                    { return nil }
func (o FakeSection) AddProperty(name string, property *kubernetes.Property) error { return nil }

func TestOutputDocumentV118(t *testing.T) {
	spec, err := kubernetes.NewSpec("../../api/v1.18/swagger.json")
	if err != nil {
		t.Errorf("Error loding swagger file")
	}

	toc, err := config.LoadTOC("../../config/v1.18/toc.yaml")
	if err != nil {
		t.Errorf("LoadTOC should not fail")
	}

	err = toc.PopulateAssociates(spec)
	if err != nil {
		t.Errorf("%s", err)
	}

	toc.OutputDocument(FakeOutput{})
}
