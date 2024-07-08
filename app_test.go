package tfmod

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestApp_Run_Dependencies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		inputs   []string
		expected string
	}{
		{
			inputs:   []string{"dependencies", "--state-dirs", "env/prd", "--base-dir=testdata/terraform", "--format=json"},
			expected: "[\"module/bar\",\"module/baz\",\"module/foo\"]\n",
		},
	}

	for _, tc := range cases {
		app := NewApp(&IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}, &Ldflags{})
		err := app.Run(tc.inputs)

		if err != nil {
			t.Fatalf("unexpected error:\n input: %v\n error: %+v", tc.inputs, err)
		}

		actual := app.IO.OutWriter.(*bytes.Buffer).String()
		if actual != tc.expected {
			t.Errorf("%s\n expected: %s actual: %s", strings.Join(tc.inputs, " "), tc.expected, actual)
		}
	}
}

func TestApp_Run_Dependents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		inputs   []string
		expected string
	}{
		{
			inputs:   []string{"dependents", "--module-dirs", "module/foo", "--base-dir=testdata/terraform", "--format=text"},
			expected: "env/dev env/prd env/stg\n",
		},
	}

	for _, tc := range cases {
		app := NewApp(&IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}, &Ldflags{})
		err := app.Run(tc.inputs)

		if err != nil {
			t.Fatalf("unexpected error:\n input: %v\n error: %+v", tc.inputs, err)
		}

		actual := app.IO.OutWriter.(*bytes.Buffer).String()
		if actual != tc.expected {
			t.Errorf("%s\n expected: %s actual: %s", strings.Join(tc.inputs, " "), tc.expected, actual)
		}
	}
}
