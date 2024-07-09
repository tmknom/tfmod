package terraform

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/testlib"
)

func TestParser_Parse(t *testing.T) {
	input := `
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

	sut := NewParser(&FakeParserStore{})
	actual, err := sut.Parse(dir.NewBaseDir("").CreateDir("env/dev"), []byte(input))
	if err != nil {
		t.Fatalf(testlib.FormatError(err, sut, input))
	}

	expected := []string{"module/bar", "module/foo"}
	for i, item := range actual {
		if item.Rel() != expected[i] {
			t.Errorf(testlib.Format(sut, expected, actual, input))
		}
	}
}

type FakeParserStore struct {
	list []string
}

func (s *FakeParserStore) Save(moduleDir *ModuleDir, tfDir *TfDir) {
	pair := strings.Join([]string{moduleDir.Rel(), tfDir.Rel()}, ":")
	s.list = append(s.list, pair)
}

func (s *FakeParserStore) GoString() string {
	return fmt.Sprintf("%v", s.list)
}
