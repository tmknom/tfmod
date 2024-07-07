package tfmod

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/tmknom/tfmod/internal/dir"
)

func TestParser_Parse(t *testing.T) {
	store := &ParserFakeStore{}
	parser := NewParser(store)
	moduleJson := `
{
  "Modules": [
    {
      "Key": "",
      "Source": "",
      "Dir": "."
    },
    {
      "Key": "foo",
      "Source": "../../module/foo",
      "Dir": "../../module/foo"
    },
    {
      "Key": "bar",
      "Source": "../bar",
      "Dir": "../../module/bar"
    }
  ]
}
`

	currentDir, _ := os.Getwd()
	baseDir := dir.NewBaseDir(currentDir)
	actual, err := parser.Parse(dir.NewDir("env/dev", baseDir), []byte(moduleJson))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"module/bar", "module/foo"}
	for i, moduleDir := range actual {
		if moduleDir.Rel() != expected[i] {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	}
}

type ParserFakeStore struct {
	list []string
}

func (s *ParserFakeStore) Actual() []string {
	sort.Strings(s.list)
	return s.list
}

func (s *ParserFakeStore) Save(moduleDir *ModuleDir, tfDir *TfDir) {
	pair := strings.Join([]string{moduleDir.Rel(), tfDir.Rel()}, ":")
	s.list = append(s.list, pair)
}

func (s *ParserFakeStore) ListTfDirs(moduleDirs []string) []string {
	return nil
}

func (s *ParserFakeStore) ListModuleDirs(stateDirs []string) []string {
	return nil
}

func (s *ParserFakeStore) Dump() {}
