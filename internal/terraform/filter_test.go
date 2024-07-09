package terraform

import (
	"testing"

	"github.com/tmknom/tfmod/internal/testlib"

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
			t.Errorf(testlib.FormatError(err, sut, tc.input))
		}

		if diff := cmp.Diff(tc.expected, dir.Dirs(actual).Slice()); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
