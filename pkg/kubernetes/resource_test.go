package kubernetes_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

func Test_ResourceLessThan(t *testing.T) {
	v1, err := kubernetes.NewAPIVersion("v1")
	if err != nil {
		t.Errorf("Error creating APIVersion from v1")
	}

	v1beta1, err := kubernetes.NewAPIVersion("v1beta1")
	if err != nil {
		t.Errorf("Error creating APIVersion from v1beta1")
	}

	v2alpha1, err := kubernetes.NewAPIVersion("v2alpha1")
	if err != nil {
		t.Errorf("Error creating APIVersion from v2alpha1")
	}

	tests := []struct {
		R1       kubernetes.Resource
		R2       kubernetes.Resource
		Expected bool
	}{
		// General case
		{
			R1: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup("apps"),
				Version: *v1,
			},
			R2: kubernetes.Resource{
				Key:     "key2",
				Group:   kubernetes.APIGroup("apps"),
				Version: *v1beta1,
			},
			Expected: true,
		},
		// Cronjob resource in v1.18
		{
			R1: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup("batch"),
				Version: *v2alpha1,
			},
			R2: kubernetes.Resource{
				Key:     "key2",
				Group:   kubernetes.APIGroup("batch"),
				Version: *v1beta1,
			},
			Expected: true,
		},
		// Event resource in v1.18
		{
			R1: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup(""),
				Version: *v1,
			},
			R2: kubernetes.Resource{
				Key:     "key2",
				Group:   kubernetes.APIGroup("events.k8s.io"),
				Version: *v1,
			},
			Expected: true,
		},
		// Ingress resource in v1.18
		{
			R1: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup("networking.k8s.io"),
				Version: *v1beta1,
			},
			R2: kubernetes.Resource{
				Key:     "key2",
				Group:   kubernetes.APIGroup("extensions"),
				Version: *v1beta1,
			},
			Expected: true,
		},
	}

	for _, test := range tests {
		result := test.R1.LessThan(&test.R2)
		if result != test.Expected {
			t.Errorf("%s < %s: expected %v but got %v", test.R1.GetGV(), test.R2.GetGV(), test.Expected, result)
		}
	}
}

func Test_ResourceGetGV(t *testing.T) {
	v1, err := kubernetes.NewAPIVersion("v1")
	if err != nil {
		t.Errorf("Error creating APIVersion from v1")
	}

	tests := []struct {
		Input    kubernetes.Resource
		Expected string
	}{
		{
			Input: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup("apps"),
				Version: *v1,
			},
			Expected: "apps/v1",
		},
		{
			Input: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup(""),
				Version: *v1,
			},
			Expected: "v1",
		},
		{
			Input: kubernetes.Resource{
				Key:     "key1",
				Group:   kubernetes.APIGroup("storage.k8s.io"),
				Version: *v1,
			},
			Expected: "storage.k8s.io/v1",
		},
	}

	for _, test := range tests {
		result := test.Input.GetGV()
		if result != test.Expected {
			t.Errorf("%#v: Expected %s but got %s", test.Input, test.Expected, result)
		}
	}
}

func Test_ResourceAdd(t *testing.T) {
	v1, err := kubernetes.NewAPIVersion("v1")
	if err != nil {
		t.Errorf("Error creating APIVersion from v1")
	}

	resources := kubernetes.ResourceMap{}
	resources.Add(&kubernetes.Resource{
		Key:     "key1",
		Group:   kubernetes.APIGroup("extensions"),
		Version: *v1,
		Kind:    kubernetes.APIKind("Kind1"),
	})
	resources.Add(&kubernetes.Resource{
		Key:     "key1",
		Group:   kubernetes.APIGroup("apps"),
		Version: *v1,
		Kind:    kubernetes.APIKind("Kind1"),
	})
	resources.Add(&kubernetes.Resource{
		Key:     "key1",
		Group:   kubernetes.APIGroup("apps"),
		Version: *v1,
		Kind:    kubernetes.APIKind("Kind2"),
	})
	if len(resources) != 2 {
		t.Errorf("Len of resources should be %d but is %d", 2, len(resources))
	}
	if _, ok := resources["Kind1"]; !ok {
		t.Errorf("Key Kind1 should exist")
	}
	if _, ok := resources["Kind2"]; !ok {
		t.Errorf("Key Kind2 should exist")
	}
	if len(resources["Kind1"]) != 2 {
		t.Errorf("Len of versions for Kind1 should be %d but is %d", 2, len(resources["Kind1"]))
	}
	if len(resources["Kind2"]) != 1 {
		t.Errorf("Len of versions for Kind2 should be %d but is %d", 1, len(resources["Kind2"]))
	}
	if resources["Kind1"][0].Group != "apps" {
		t.Errorf("Recent version for Kind1 should be %s but is %s", "apps", resources["Kind1"][0].Group)
	}
	if resources["Kind1"][1].Group != "extensions" {
		t.Errorf("Previous version for Kind1 should be %s but is %s", "extensions", resources["Kind1"][1].Group)
	}
}