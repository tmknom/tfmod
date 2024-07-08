package terraform

import (
	"testing"

	"github.com/tmknom/tfmod/internal/dir"
)

func TestFilter_SubDirs(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "../../testdata/terraform/env",
			expected: []string{"dev", "prd", "stg"},
		},
		{
			input:    "../../testdata/terraform",
			expected: []string{"env/dev", "env/prd", "env/stg", "module/bar", "module/baz", "module/foo"},
		},
	}

	for _, tc := range cases {
		baseDir := dir.NewBaseDir(tc.input)
		sut := NewFilter(baseDir)

		actual, err := sut.SubDirs()
		if err != nil {
			t.Errorf("unexpected error:\n input: %v\n error: %+v", tc.input, err)
		}

		if len(tc.expected) != len(actual) {
			t.Fatalf("error length:\n input: %v\n expected: %v, actual: %v", tc.input, tc.expected, actual)
		}
		for i, item := range tc.expected {
			if item != actual[i].Rel() {
				t.Errorf("error index: \n input: %v\n index: %d, expected: %v, actual: %v", tc.input, i, tc.expected, actual)
			}
		}
	}
}
