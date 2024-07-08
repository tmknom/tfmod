package dir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tmknom/tfmod/internal/testlib"
)

func TestDir_Abs(t *testing.T) {
	currentDir, _ := os.Getwd()
	type input struct {
		raw  string
		base string
	}

	cases := []struct {
		input    input
		expected string
	}{
		{
			input: input{
				raw:  ".",
				base: currentDir,
			},
			expected: currentDir,
		},
		{
			input: input{
				raw:  filepath.Join(currentDir, "path/to/dir"),
				base: currentDir,
			},
			expected: filepath.Join(currentDir, "path/to/dir"),
		},
		{
			input: input{
				raw:  "path/to/dir",
				base: currentDir,
			},
			expected: filepath.Join(currentDir, "path/to/dir"),
		},
		{
			input: input{
				raw:  "./path/to/dir/../../test",
				base: currentDir,
			},
			expected: filepath.Join(currentDir, "path/test"),
		},
		{
			input: input{
				raw:  "path/to/dir",
				base: "/test/base",
			},
			expected: "/test/base/path/to/dir",
		},
		{
			input: input{
				raw:  "path/to/dir",
				base: "test/base",
			},
			expected: filepath.Join(currentDir, "test/base/path/to/dir"),
		},
		{
			input: input{
				raw:  "test/base/path/to/dir",
				base: "test/base",
			},
			expected: filepath.Join(currentDir, "test/base/path/to/dir"),
		},
		{
			input: input{
				raw:  "test/base/../base/path/to/dir",
				base: "test/base",
			},
			expected: filepath.Join(currentDir, "test/base/path/to/dir"),
		},
	}

	for _, tc := range cases {
		sut := NewDir(tc.input.raw, NewBaseDir(tc.input.base))
		actual := sut.Abs()
		if actual != tc.expected {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func TestDir_Rel(t *testing.T) {
	currentDir, _ := os.Getwd()

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    currentDir,
			expected: ".",
		},
		{
			input:    ".",
			expected: ".",
		},
		{
			input:    "path/to/dir",
			expected: "path/to/dir",
		},
		{
			input:    "./path/to/dir/../../test",
			expected: "path/test",
		},
		{
			input:    filepath.Join(currentDir, "path/to/dir"),
			expected: "path/to/dir",
		},
	}

	baseDir := NewBaseDir(currentDir)
	for _, tc := range cases {
		sut := NewDir(tc.input, baseDir)
		actual := sut.Rel()
		if actual != tc.expected {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
