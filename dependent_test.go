package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/terraform"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestDependentRunner_List(t *testing.T) {
	cases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"module/bar"},
			expected: []string{"env/prd"},
		},
		{
			input:    []string{"module/baz"},
			expected: []string{"env/prd"},
		},
		{
			input:    []string{"module/bar", "module/baz"},
			expected: []string{"env/prd"},
		},
		{
			input:    []string{"module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
		{
			input:    []string{"module/bar", "module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
	}

	globalFlags := &GlobalFlags{
		base: "testdata/terraform",
	}

	for _, tc := range cases {
		bufIO := &IO{InReader: &bytes.Buffer{}, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &DependentFlags{
			ModulePaths: tc.input,
			GlobalFlags: globalFlags,
		}
		sut := NewDependentRunner(flags, terraform.NewDependentStore(), bufIO)

		actual, err := sut.List()
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
