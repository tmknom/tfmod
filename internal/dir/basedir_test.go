package dir

import (
	"fmt"
	"os"
	"testing"

	"github.com/tmknom/tfmod/internal/testlib"

	"github.com/google/go-cmp/cmp"
)

func TestBaseDir_Abs(t *testing.T) {
	currentDir, _ := os.Getwd()

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    ".",
			expected: currentDir,
		},
		{
			input:    "testdata/terraform/../",
			expected: fmt.Sprintf("%s/%s", currentDir, "testdata"),
		},
		{
			input:    "/path/to/dir",
			expected: "/path/to/dir",
		},
		{
			input:    "../../internal/dir/foo/bar/baz/../",
			expected: fmt.Sprintf("%s/%s", currentDir, "foo/bar"),
		},
	}

	for _, tc := range cases {
		sut := NewBaseDir(tc.input)
		actual := sut.Abs()
		if actual != tc.expected {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func TestBaseDir_FilterSubDirs(t *testing.T) {
	rel := "../../testdata/terraform/env/"
	sut := NewBaseDir(rel)

	actual, err := sut.FilterSubDirs(".tf", ".terraform")
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	expected := []string{"dev", "prd", "stg"}
	if diff := cmp.Diff(expected, Dirs(actual).Slice()); diff != "" {
		t.Errorf(testlib.FormatWithoutInput(sut, expected, actual))
	}
}
