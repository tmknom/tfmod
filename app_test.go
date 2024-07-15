package tfmod

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/tmknom/tfmod/internal/testlib"
)

func TestApp_Run(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"download", "--base=testdata/terraform/env", "--max-concurrent=4"},
			expected: "dev prd stg\n",
		},
		{
			input:    []string{"dependency", "--state", "env/prd", "--base=testdata/terraform", "--format=json"},
			expected: "[\"module/bar\",\"module/baz\",\"module/foo\"]\n",
		},
		{
			input:    []string{"dependent", "--module", "module/foo", "--base=testdata/terraform", "--format=text"},
			expected: "env/dev env/prd env/stg\n",
		},
	}

	for _, tc := range cases {
		sut := NewApp(&IO{InReader: os.Stdin, OutWriter: &bytes.Buffer{}, ErrWriter: os.Stderr}, &Ldflags{})
		err := sut.Run(context.Background(), tc.input)

		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		actual := sut.IO.OutWriter.(*bytes.Buffer).String()
		if actual != tc.expected {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
