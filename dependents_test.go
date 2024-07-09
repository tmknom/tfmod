package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/terraform"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestDependentsRunner_List(t *testing.T) {
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
		BaseDir: "testdata/terraform",
		Debug:   true,
	}

	for _, tc := range cases {
		bufIO := &IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &DependentsFlags{
			ModuleDirs:  tc.input,
			GlobalFlags: globalFlags,
		}
		sut := NewDependentsRunner(flags, terraform.NewDependentStore(), bufIO)

		actual, err := sut.List()
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
