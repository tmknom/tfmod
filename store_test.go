package tfmod

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/dir"
)

func TestInMemoryStore_ListModuleDirs(t *testing.T) {
	cases := []struct {
		inputs   []string
		expected []string
	}{
		{
			inputs:   []string{"env/dev"},
			expected: []string{"module/bar", "module/foo"},
		},
		{
			inputs:   []string{"env/dev", "env/prd"},
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
		actual := sut.ListModuleDirs(baseDir.ConvertDirs(tc.inputs))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("input: %v, expected: %v, actual: %v", tc.inputs, tc.expected, actual)
		}
	}
}

func TestInMemoryStore_ListTfDirs(t *testing.T) {
	cases := []struct {
		inputs   []string
		expected []string
	}{
		{
			inputs:   []string{"module/foo"},
			expected: []string{"env/dev", "env/prd", "env/stg"},
		},
		{
			inputs:   []string{"module/bar", "module/baz"},
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
		actual := sut.ListTfDirs(baseDir.ConvertDirs(tc.inputs))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("input: %v, expected: %v, actual: %v", tc.inputs, tc.expected, actual)
		}
	}
}
