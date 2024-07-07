package tfmod

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDependenciesRunner_List(t *testing.T) {
	cases := []struct {
		dirs     []string
		expected []string
	}{
		{
			dirs:     []string{"env/dev"},
			expected: []string{"module/foo"},
		},
		{
			dirs:     []string{"env/dev", "env/stg"},
			expected: []string{"module/foo"},
		},
		{
			dirs:     []string{"env/prd"},
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
			StateDirs:   tc.dirs,
			GlobalFlags: globalFlags,
		}
		runner := NewDependenciesRunner(flags, NewInMemoryStore(), bufIO)

		actual, err := runner.List()
		if err != nil {
			t.Fatalf("%s: unexpected error: %e", strings.Join(tc.dirs, " "), err)
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("dirs: %v, expected: %v, actual: %v", tc.dirs, tc.expected, actual)
		}
	}
}
