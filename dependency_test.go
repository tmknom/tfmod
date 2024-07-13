package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/terraform"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestDependencyRunner_List(t *testing.T) {
	cases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"env/dev"},
			expected: []string{"module/foo"},
		},
		{
			input:    []string{"env/dev", "env/stg"},
			expected: []string{"module/foo"},
		},
		{
			input:    []string{"env/prd"},
			expected: []string{"module/bar", "module/baz", "module/foo"},
		},
	}

	globalFlags := &GlobalFlags{
		base: "testdata/terraform",
	}

	for _, tc := range cases {
		bufIO := &IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &DependencyFlags{
			StatePaths:  tc.input,
			GlobalFlags: globalFlags,
		}
		sut := NewDependencyRunner(flags, terraform.NewDependencyStore(), bufIO)

		actual, err := sut.List()
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
