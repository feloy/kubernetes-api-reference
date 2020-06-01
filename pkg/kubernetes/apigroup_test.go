package kubernetes_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

func Test_GoImportPrefix(t *testing.T) {
	tests := []struct {
		Group    kubernetes.APIGroup
		Expected string
	}{
		{
			Group:    "io.k8s.api.core",
			Expected: "k8s.io/api/core",
		},
	}

	for _, test := range tests {
		result := test.Group.GoImportPrefix()
		if result != test.Expected {
			t.Errorf("%s: Expected %s but got %s", test.Group, test.Expected, result)
		}
	}
}

func Test_APIGroupReplaces(t *testing.T) {
	tests := []struct {
		Group1   kubernetes.APIGroup
		Group2   kubernetes.APIGroup
		Expected bool
	}{
		{"policy", "extensions", true},
		{"admissionregistration.k8s.io", "apiextensions.k8s.io", true},
	}

	for _, test := range tests {
		result := test.Group1.Replaces(test.Group2)
		if result != test.Expected {
			t.Errorf("%s replaces %s: expected %v but got %v", test.Group1, test.Group2, test.Expected, result)
		}
	}
}
