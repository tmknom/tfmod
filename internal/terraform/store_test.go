package terraform

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestDependencyStore_List(t *testing.T) {
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
	sut := NewDependencyStore()
	sut.Save(NewModuleDir("module/foo", baseDir), NewStateDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/foo", baseDir), NewStateDir("env/prd", baseDir))
	sut.Save(NewModuleDir("module/bar", baseDir), NewStateDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/baz", baseDir), NewStateDir("env/prd", baseDir))

	for _, tc := range cases {
		actual := sut.List(FactoryTestDirs(tc.input, baseDir))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func TestDependentStore_List(t *testing.T) {
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
	sut := NewDependentStore()
	sut.Save(NewModuleDir("module/foo", baseDir), NewStateDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/foo", baseDir), NewStateDir("env/prd", baseDir))
	sut.Save(NewModuleDir("module/foo", baseDir), NewStateDir("env/stg", baseDir))
	sut.Save(NewModuleDir("module/bar", baseDir), NewStateDir("env/dev", baseDir))
	sut.Save(NewModuleDir("module/baz", baseDir), NewStateDir("env/prd", baseDir))

	for _, tc := range cases {
		actual := sut.List(FactoryTestDirs(tc.input, baseDir))

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func FactoryTestDirs(paths []string, baseDir *dir.BaseDir) []*dir.Dir {
	result := make([]*dir.Dir, 0, len(paths))
	for _, path := range paths {
		result = append(result, baseDir.CreateDir(path))
	}
	return result
}
