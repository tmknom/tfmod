package tfmod

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDependenciesRunner_List(t *testing.T) {
	cases := []struct {
		inputs   []string
		expected []string
	}{
		{
			inputs:   []string{"env/dev"},
			expected: []string{"module/foo"},
		},
		{
			inputs:   []string{"env/dev", "env/stg"},
			expected: []string{"module/foo"},
		},
		{
			inputs:   []string{"env/prd"},
			expected: []string{"module/bar", "module/baz", "module/foo"},
		},
	}

	globalFlags := &GlobalFlags{
		BaseDir:  "testdata/terraform",
		EnableTf: false,
		Debug:    true,
	}

	for _, tc := range cases {
		bufIO := &IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}
		flags := &DependenciesFlags{
			StateDirs:   tc.inputs,
			GlobalFlags: globalFlags,
		}
		runner := NewDependenciesRunner(flags, NewInMemoryStore(), bufIO)

		actual, err := runner.List()
		if err != nil {
			t.Fatalf("unexpected error:\n input: %v\n error: %+v", tc.inputs, err)
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("input: %v, expected: %v, actual: %v", tc.inputs, tc.expected, actual)
		}
	}
}
