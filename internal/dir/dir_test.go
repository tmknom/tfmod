package dir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDir_Abs(t *testing.T) {
	currentDir, _ := os.Getwd()

	cases := []struct {
		raw      string
		expected string
	}{
		{
			raw:      currentDir,
			expected: currentDir,
		},
		{
			raw:      ".",
			expected: currentDir,
		},
		{
			raw:      "path/to/dir",
			expected: filepath.Join(currentDir, "path/to/dir"),
		},
		{
			raw:      "./path/to/dir/../../test",
			expected: filepath.Join(currentDir, "path/test"),
		},
		{
			raw:      filepath.Join(currentDir, "path/to/dir"),
			expected: filepath.Join(currentDir, "path/to/dir"),
		},
	}

	baseDir := NewBaseDir(currentDir)
	for _, tc := range cases {
		sut := NewDir(tc.raw, baseDir)
		actual := sut.Abs()
		if actual != tc.expected {
			t.Errorf("raw: %s, expected: %v, actual: %v", tc.raw, tc.expected, actual)
		}
	}
}

func TestDir_Rel(t *testing.T) {
	currentDir, _ := os.Getwd()

	cases := []struct {
		raw      string
		expected string
	}{
		{
			raw:      currentDir,
			expected: ".",
		},
		{
			raw:      ".",
			expected: ".",
		},
		{
			raw:      "path/to/dir",
			expected: "path/to/dir",
		},
		{
			raw:      "./path/to/dir/../../test",
			expected: "path/test",
		},
		{
			raw:      filepath.Join(currentDir, "path/to/dir"),
			expected: "path/to/dir",
		},
	}

	baseDir := NewBaseDir(currentDir)
	for _, tc := range cases {
		sut := NewDir(tc.raw, baseDir)
		actual := sut.Rel()
		if actual != tc.expected {
			t.Errorf("raw: %s, expected: %v, actual: %v", tc.raw, tc.expected, actual)
		}
	}
}
