package kubernetes_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

func Test_SpecV118(t *testing.T) {
	spec, err := kubernetes.NewSpec("../../api/v1.18/swagger.json")
	if err != nil {
		t.Errorf("NewSpec should not return an errors but returns %s", err)
	}
	if len(*spec.Resources) != 114 {
		t.Errorf("Spec should contain %d resources but contains %d", 114, len(*spec.Resources))
	}
}
