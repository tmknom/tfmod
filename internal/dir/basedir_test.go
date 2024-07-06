package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	expected := []string{
		filepath.Join(sut.Abs(), "dev"),
		filepath.Join(sut.Abs(), "prd"),
		filepath.Join(sut.Abs(), "stg"),
	}

	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
