package config_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/config"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

func Test_LoadTOCv118(t *testing.T) {
	toc, err := config.LoadTOC("../../config/v1.18/toc.yaml")
	if err != nil {
		t.Errorf("should not get an error but got: %s", err)
	}
	if len(toc.Parts) != 8 {
		t.Errorf("Should get %d parts but got %d", 8, len(toc.Parts))
	}
}

func Test_PopulateAssociatesv118(t *testing.T) {
	spec, err := kubernetes.NewSpec("../../api/v1.18/swagger.json")
	if err != nil {
		t.Errorf("Error loding swagger file")
	}

	toc, err := config.LoadTOC("../../config/v1.18/toc.yaml")
	if err != nil {
		t.Errorf("Error loding toc file")
	}

	toc.PopulateAssociates(spec)

	l := len(toc.Parts[0].Chapters[0].Sections)
	if l != 4 {
		t.Errorf("Pod chapter should contain %d sections but contains %d sections", 4, l)
	}
}
