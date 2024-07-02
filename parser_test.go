package tfmod

import (
	"os"
	"sort"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	store := &ParserFakeStore{}
	currentDir, _ := os.Getwd()
	parser := NewParser(BaseDir(currentDir), store)
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

	err := parser.Parse("env/dev", []byte(moduleJson))
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	expected := "module/bar:env/dev,module/foo:env/dev"
	if store.Actual() != expected {
		t.Errorf("expected: %s, actual: %s", expected, store.Actual())
	}
}

type ParserFakeStore struct {
	list []string
}

func (s *ParserFakeStore) Actual() string {
	sort.Strings(s.list)
	return strings.Join(s.list, ",")
}

func (s *ParserFakeStore) Save(moduleDir ModuleDir, tfDir TfDir) {
	pair := strings.Join([]string{moduleDir, tfDir}, ":")
	s.list = append(s.list, pair)
}

func (s *ParserFakeStore) List(sourceDirs SourceDirs) *TfDirs {
	return nil
}

func (s *ParserFakeStore) Dump() *DependentMap {
	return nil
}
