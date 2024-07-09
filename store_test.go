package tfmod

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/terraform"
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
	sut.Save(terraform.NewModuleDir("module/foo", baseDir), terraform.NewTfDir("env/dev", baseDir))
	sut.Save(terraform.NewModuleDir("module/foo", baseDir), terraform.NewTfDir("env/prd", baseDir))
	sut.Save(terraform.NewModuleDir("module/bar", baseDir), terraform.NewTfDir("env/dev", baseDir))
	sut.Save(terraform.NewModuleDir("module/baz", baseDir), terraform.NewTfDir("env/prd", baseDir))

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
	sut.Save(terraform.NewModuleDir("module/foo", baseDir), terraform.NewTfDir("env/dev", baseDir))
	sut.Save(terraform.NewModuleDir("module/foo", baseDir), terraform.NewTfDir("env/prd", baseDir))
	sut.Save(terraform.NewModuleDir("module/foo", baseDir), terraform.NewTfDir("env/stg", baseDir))
	sut.Save(terraform.NewModuleDir("module/bar", baseDir), terraform.NewTfDir("env/dev", baseDir))
	sut.Save(terraform.NewModuleDir("module/baz", baseDir), terraform.NewTfDir("env/prd", baseDir))

	for _, tc := range cases {
		actual := sut.ListTfDirs(baseDir.ConvertDirs(tc.input))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
