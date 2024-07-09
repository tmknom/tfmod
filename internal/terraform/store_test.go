package terraform

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestInMemoryStore_ListModuleDirs(t *testing.T) {
	cases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"env/dev"},
			expected: []string{"module/bar", "module/foo"},
		},
		{
			input:    []string{"env/dev", "env/prd"},
			expected: []string{"module/bar", "module/baz", "module/foo"},
		},
	}

	baseDir := dir.NewBaseDir("testdata/terraform")
	sut := NewInMemoryStore()
	sut.Save(NewModuleDir("module/foo", baseDir), NewTfDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/foo", baseDir), NewTfDir("env/prd", baseDir))
	sut.Save(NewModuleDir("module/bar", baseDir), NewTfDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/baz", baseDir), NewTfDir("env/prd", baseDir))

	for _, tc := range cases {
		actual := sut.ListModuleDirs(baseDir.ConvertDirs(tc.input))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func TestInMemoryStore_ListTfDirs(t *testing.T) {
	cases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
		{
			input:    []string{"module/bar", "module/baz"},
			expected: []string{"env/dev", "env/prd"},
		},
	}

	baseDir := dir.NewBaseDir("testdata/terraform")
	sut := NewInMemoryStore()
	sut.Save(NewModuleDir("module/foo", baseDir), NewTfDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/foo", baseDir), NewTfDir("env/prd", baseDir))
	sut.Save(NewModuleDir("module/foo", baseDir), NewTfDir("env/stg", baseDir))
	sut.Save(NewModuleDir("module/bar", baseDir), NewTfDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/baz", baseDir), NewTfDir("env/prd", baseDir))

	for _, tc := range cases {
		actual := sut.ListTfDirs(baseDir.ConvertDirs(tc.input))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
