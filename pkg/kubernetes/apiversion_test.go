package kubernetes_test

import (
	"reflect"
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

var (
	two   = 2
	three = 3
)

func Test_NewAPIVersion(t *testing.T) {
	tests := []struct {
		Input         string
		Expected      *kubernetes.APIVersion
		ExpectedInErr bool
	}{
		{
			Input: "v1",
			Expected: &kubernetes.APIVersion{
				Version:      1,
				Stage:        kubernetes.StageGA,
				StageVersion: nil,
			},
			ExpectedInErr: false,
		},
		{
			Input: "v1alpha2",
			Expected: &kubernetes.APIVersion{
				Version:      1,
				Stage:        kubernetes.StageAlpha,
				StageVersion: &two,
			},
			ExpectedInErr: false,
		},
		{
			Input: "v2beta3",
			Expected: &kubernetes.APIVersion{
				Version:      2,
				Stage:        kubernetes.StageBeta,
				StageVersion: &three,
			},
			ExpectedInErr: false,
		},
		{
			Input:         "_v1_",
			Expected:      nil,
			ExpectedInErr: true,
		},
		{
			Input:         "v1alpha",
			Expected:      nil,
			ExpectedInErr: true,
		},
		{
			Input:         "va",
			Expected:      nil,
			ExpectedInErr: true,
		},
		{
			Input:         "v1gamma2",
			Expected:      nil,
			ExpectedInErr: true,
		},
		{
			Input:         "v1gamma",
			Expected:      nil,
			ExpectedInErr: true,
		},
	}

	for _, test := range tests {
		result, err := kubernetes.NewAPIVersion(test.Input)
		if (err != nil) != test.ExpectedInErr {
			t.Errorf("%s: Expected error %v but got %v", test.Input, test.ExpectedInErr, err != nil)
		}
		if !reflect.DeepEqual(result, test.Expected) {
			t.Errorf("%s: Expected result is %v but got %v", test.Input, test.Expected, result)
		}
	}
}

func Test_String(t *testing.T) {
	tests := []struct {
		Input    *kubernetes.APIVersion
		Expected string
	}{
		{
			Input: &kubernetes.APIVersion{
				Version: 1,
				Stage:   kubernetes.StageGA,
			},
			Expected: "v1",
		},
		{
			Input: &kubernetes.APIVersion{
				Version:      1,
				Stage:        kubernetes.StageAlpha,
				StageVersion: &two,
			},
			Expected: "v1alpha2",
		},
		{
			Input:    nil,
			Expected: "",
		},
	}

	for _, test := range tests {
		result := test.Input.String()
		if result != test.Expected {
			t.Errorf("%#v: Expected %s but got %s", test.Input, test.Expected, result)
		}
	}
}

func Test_LessThan(t *testing.T) {
	tests := []struct {
		V1       string
		V2       string
		Expected bool
	}{
		{"v1", "v2", true},
		{"v1", "v1alpha3", false},
		{"v1", "v1beta2", false},
		{"v1", "v1", false},
		{"v1", "v2alpha1", true},
	}

	for _, test := range tests {
		v1, err := kubernetes.NewAPIVersion(test.V1)
		if err != nil {
			t.Errorf("Creating an APIVersion with '%s' should work", test.V1)
		}
		v2, err := kubernetes.NewAPIVersion(test.V2)
		if err != nil {
			t.Errorf("Creating an APIVersion with '%s' should work", test.V2)
		}
		result := v1.LessThan(v2)
		if result != test.Expected {
			t.Errorf("%s < %s: Expected %v but got %v", test.V1, test.V2, test.Expected, result)
		}
	}
}

func Test_Replaces(t *testing.T) {
	tests := []struct {
		V1       string
		V2       string
		Expected bool
	}{
		{"v2", "v1", false},
		{"v1", "v2", false},
		{"v1", "v1alpha3", true},
		{"v1", "v1beta2", true},
		{"v1beta1", "v1alpha3", true},
		{"v1beta2", "v1beta1", true},
		{"v1", "v1", false},
		{"v1", "v2alpha1", false},
		{"v2alpha1", "v1", false},
	}

	for _, test := range tests {
		v1, err := kubernetes.NewAPIVersion(test.V1)
		if err != nil {
			t.Errorf("Creating an APIVersion with '%s' should work", test.V1)
		}
		v2, err := kubernetes.NewAPIVersion(test.V2)
		if err != nil {
			t.Errorf("Creating an APIVersion with '%s' should work", test.V2)
		}
		result := v1.Replaces(v2)
		if result != test.Expected {
			t.Errorf("%s replaces %s: Expected %v but got %v", test.V1, test.V2, test.Expected, result)
		}
	}
}