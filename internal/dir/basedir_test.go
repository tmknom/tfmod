package dir

import (
	"fmt"
	"os"
	"testing"
)

func TestBaseDir_Abs(t *testing.T) {
	currentDir, _ := os.Getwd()

	cases := []struct {
		raw      string
		expected string
	}{
		{
			raw:      ".",
			expected: currentDir,
		},
		{
			raw:      "testdata/terraform/../",
			expected: fmt.Sprintf("%s/%s", currentDir, "testdata"),
		},
		{
			raw:      "/path/to/dir",
			expected: "/path/to/dir",
		},
	}

	for _, tc := range cases {
		sut := NewBaseDir(tc.raw)
		actual := sut.Abs()
		if actual != tc.expected {
			t.Errorf("expected: %v, actual: %v", tc.expected, actual)
		}
	}
}

func TestBaseDir_FilterSubDirs(t *testing.T) {
	rel := "../../testdata/terraform/env/"
	sut := NewBaseDir(rel)

	actual, err := sut.FilterSubDirs(".tf", ".terraform")
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	expected := []string{"dev", "prd", "stg"}
	if len(expected) != len(actual) {
		t.Fatalf("diff length, expected: %v, actual: %v", expected, actual)
	}
	for i, item := range expected {
		if item != actual[i].Rel() {
			t.Errorf("diff index: %d, expected: %v, actual: %v", i, expected, actual)
		}
	}
}
