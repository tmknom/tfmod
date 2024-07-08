package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDependentsRunner_List(t *testing.T) {
	cases := []struct {
		inputs   []string
		expected []string
	}{
		{
			inputs:   []string{"module/bar"},
			expected: []string{"env/prd"},
		},
		{
			inputs:   []string{"module/baz"},
			expected: []string{"env/prd"},
		},
		{
			inputs:   []string{"module/bar", "module/baz"},
			expected: []string{"env/prd"},
		},
		{
			inputs:   []string{"module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
		{
			inputs:   []string{"module/bar", "module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
	}

	globalFlags := &GlobalFlags{
		BaseDir:  "testdata/terraform",
		EnableTf: false,
		Debug:    true,
	}

	for _, tc := range cases {
		bufIO := &IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &DependentsFlags{
			ModuleDirs:  tc.inputs,
			GlobalFlags: globalFlags,
		}
		runner := NewDependentsRunner(flags, NewInMemoryStore(), bufIO)

		actual, err := runner.List()
		if err != nil {
			t.Fatalf("unexpected error:\n input: %v\n error: %+v", tc.inputs, err)
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("input: %v, expected: %v, actual: %v", tc.inputs, tc.expected, actual)
		}
	}
}
