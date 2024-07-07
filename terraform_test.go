package tfmod

import (
	"testing"

	"github.com/tmknom/tfmod/internal/dir"
)

func TestTerraform_GetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		baseDir  string
		expected []string
	}{
		{
			baseDir:  "testdata/terraform/env",
			expected: []string{"dev", "prd", "stg"},
		},
	}

	for _, tc := range cases {
		sut := NewTerraform(dir.NewBaseDir(tc.baseDir), true)

		actual, err := sut.GetAll()
		if err != nil {
			t.Fatalf("%v: unexpected error: %e", tc.baseDir, err)
		}

		if len(tc.expected) != len(actual) {
			t.Fatalf("diff length, expected: %v, actual: %v", tc.expected, actual)
		}
		for i, item := range tc.expected {
			if item != actual[i].Rel() {
				t.Errorf("diff index: %d, expected: %v, actual: %v", i, tc.expected, actual)
			}
		}
	}
}
