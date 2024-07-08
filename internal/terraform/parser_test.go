package terraform

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/tmknom/tfmod/internal/dir"
)

func TestParser_Parse(t *testing.T) {
	store := &FakeParserStore{}
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
	actual, err := parser.Parse(baseDir.CreateDir("env/dev"), []byte(moduleJson))
	if err != nil {
		t.Fatalf("unexpected error:\n input: %v\n error: %+v", moduleJson, err)
	}

	expected := []string{"module/bar", "module/foo"}
	for i, moduleDir := range actual {
		if moduleDir.Rel() != expected[i] {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	}
}

type FakeParserStore struct {
	list []string
}

func (s *FakeParserStore) Actual() []string {
	sort.Strings(s.list)
	return s.list
}

func (s *FakeParserStore) Save(moduleDir *ModuleDir, tfDir *TfDir) {
	pair := strings.Join([]string{moduleDir.Rel(), tfDir.Rel()}, ":")
	s.list = append(s.list, pair)
}
