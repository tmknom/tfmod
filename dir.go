package tfmod

import (
	"encoding/json"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

type BaseDir string

func (d *BaseDir) String() string {
	return string(*d)
}

func (d *BaseDir) ListTfDirs() (TfDirs, error) {
	tfDirs := NewTfDirs()
	err := filepath.WalkDir(d.String(), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".tf" {
			tfDirs.Add(filepath.Dir(path))
		}
		return nil
	})
	return tfDirs, err
}

type TfDir = string
type ModuleDir = string
type SourceDir = string

type SourceDirs []SourceDir

func NewSourceDirs(inputs []string) SourceDirs {
	return inputs
}

type TfDirs map[TfDir]bool

func NewTfDirs() TfDirs {
	return make(TfDirs, 64)
}

func (d *TfDirs) Add(tfDir TfDir) {
	(*d)[tfDir] = true
}

func (d *TfDirs) SortedSlice() []string {
	result := make([]string, 0, len(*d))
	for tfDir := range *d {
		result = append(result, tfDir)
	}
	sort.Strings(result)
	return result
}

func (d *TfDirs) String() string {
	return strings.Join(d.SortedSlice(), " ")
}

func (d *TfDirs) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.SortedSlice())
}

func (d *TfDirs) ToJson() string {
	return SimpleJsonMarshal(d)
}

type DependentMap map[ModuleDir][]TfDir

func NewDependentMap() DependentMap {
	return make(DependentMap, 64)
}

func (m *DependentMap) IsModule(dir string) bool {
	_, ok := (*m)[dir]
	return ok
}

func (m *DependentMap) ToJson() string {
	return SimpleJsonMarshal(m)
}

func SimpleJsonMarshal(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatalln(err)
	}
	return string(bytes)
}
