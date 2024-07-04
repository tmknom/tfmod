package tfmod

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser_Parse(t *testing.T) {
	store := &ParserFakeStore{}
	currentDir, _ := os.Getwd()
	parser := NewParser(NewBaseDir(currentDir), store)
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

	actual, err := parser.Parse("env/dev", []byte(moduleJson))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"module/bar", "module/foo"}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

type ParserFakeStore struct {
	list []string
}

func (s *ParserFakeStore) Actual() []string {
	sort.Strings(s.list)
	return s.list
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
