package tfmod

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestAppRunWithDependencies(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"dependencies", "--state-dirs", "testdata/terraform/env/dev"},
			expected: "[\"testdata/terraform/module/foo\"]\n",
		},
	}

	for _, tc := range cases {
		app := NewApp(&IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}, &Ldflags{})
		err := app.Run(tc.args)

		if err != nil {
			t.Fatalf("%s: unexpected error: %e", strings.Join(tc.args, " "), err)
		}

		actual := app.IO.OutWriter.(*bytes.Buffer).String()
		if actual != tc.expected {
			t.Errorf("%s: expected: %s, actual: %s", strings.Join(tc.args, " "), tc.expected, actual)
		}
	}
}

func TestAppRunWithDependents(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"dependents", "--module-dirs", "testdata/terraform/module/foo"},
			expected: "[\"testdata/terraform/env/dev\",\"testdata/terraform/env/prd\",\"testdata/terraform/env/stg\"]\n",
		},
		{
			args:     []string{"dependents", "--module-dirs", "testdata/terraform/module/bar,testdata/terraform/module/baz"},
			expected: "[\"testdata/terraform/env/prd\"]\n",
		},
	}

	for _, tc := range cases {
		app := NewApp(&IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}, &Ldflags{})
		err := app.Run(tc.args)

		if err != nil {
			t.Fatalf("%s: unexpected error: %e", strings.Join(tc.args, " "), err)
		}

		actual := app.IO.OutWriter.(*bytes.Buffer).String()
		if actual != tc.expected {
			t.Errorf("%s: expected: %s, actual: %s", strings.Join(tc.args, " "), tc.expected, actual)
		}
	}
}
