package tfmod

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestApp_Run_Dependencies(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"dependencies", "--state-dirs", "env/prd", "--base-dir=testdata/terraform", "--enable-tf=false", "--format=json"},
			expected: "[\"module/bar\",\"module/baz\",\"module/foo\"]\n",
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
			t.Errorf("%s\n expected: %s actual: %s", strings.Join(tc.args, " "), tc.expected, actual)
		}
	}
}

func TestApp_Run_Dependents(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"dependents", "--module-dirs", "module/foo", "--base-dir=testdata/terraform", "--enable-tf=false", "--format=text"},
			expected: "env/dev env/prd env/stg\n",
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
			t.Errorf("%s\n expected: %s actual: %s", strings.Join(tc.args, " "), tc.expected, actual)
		}
	}
}
