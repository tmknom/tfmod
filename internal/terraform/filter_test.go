package terraform

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

		if diff := cmp.Diff(tc.expected, dir.Dirs(actual).Slice()); diff != "" {
			t.Errorf("expected: %v, actual: %v", tc.expected, actual)
		}
	}
}
