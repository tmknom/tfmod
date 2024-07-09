package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestGetRunner_List(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "testdata/terraform/env",
			expected: []string{"dev", "prd", "stg"},
		},
	}

	for _, tc := range cases {
		bufIO := &IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &GetFlags{
			GlobalFlags: &GlobalFlags{
				BaseDir: tc.input,
				Debug:   true,
			},
		}
		sut := NewGetRunner(flags, bufIO)

		actual, err := sut.TerraformGet()
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
