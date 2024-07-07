package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetRunner_List(t *testing.T) {
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
		bufIO := &IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &GetFlags{
			GlobalFlags: &GlobalFlags{
				BaseDir:  tc.baseDir,
				EnableTf: false,
				Debug:    true,
			},
		}
		runner := NewGetRunner(flags, NewInMemoryStore(), bufIO)

		actual, err := runner.TerraformGet()
		if err != nil {
			t.Fatalf("%v: unexpected error: %e", tc.baseDir, err)
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("flags: %#v, expected: %v, actual: %v", flags.GlobalFlags, tc.expected, actual)
		}
	}
}
