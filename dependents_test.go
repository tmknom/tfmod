package tfmod

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDependentsRunner_List(t *testing.T) {
	cases := []struct {
		dirs     []string
		expected []string
	}{
		{
			dirs:     []string{"module/bar"},
			expected: []string{"env/prd"},
		},
		{
			dirs:     []string{"module/baz"},
			expected: []string{"env/prd"},
		},
		{
			dirs:     []string{"module/bar", "module/baz"},
			expected: []string{"env/prd"},
		},
		{
			dirs:     []string{"module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
		{
			dirs:     []string{"module/bar", "module/foo"},
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
			ModuleDirs:  tc.dirs,
			GlobalFlags: globalFlags,
		}
		runner := NewDependentsRunner(flags, NewInMemoryStore(), bufIO)

		actual, err := runner.List()
		if err != nil {
			t.Fatalf("%s: unexpected error: %e", strings.Join(tc.dirs, " "), err)
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("dirs: %v, expected: %v, actual: %v", tc.dirs, tc.expected, actual)
		}
	}
}
