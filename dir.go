package tfmod

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tmknom/tfmod/internal/errlib"
)

type BaseDir struct {
	Raw string
}

func NewBaseDir(raw string) *BaseDir {
	return &BaseDir{
		Raw: raw,
	}
}

func (d *BaseDir) String() string {
	return d.Abs()
}

func (d *BaseDir) Abs() string {
	dir := d.Raw
	if len(dir) > 0 && dir[0] != '/' {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalln("invalid current directory:", err)
		}
		dir = filepath.Join(currentDir, dir)
	}

	result, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln("invalid base dir:", err)
	}
	return result
}

func (d *BaseDir) ListTfDirs() (*TfDirs, error) {
	tfDirs := NewTfDirs()
	err := filepath.WalkDir(d.Abs(), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return errlib.Wrapf(err, "invalid base dir: %#v", d)
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".tf" {
			relPath, err := filepath.Rel(d.Abs(), path)
			if err != nil {
				return errlib.Wrapf(err, "invalid base dir: %#v", d)
			}
			tfDirs.Add(filepath.Dir(relPath))
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

type TfDirs struct {
	set  map[TfDir]bool
	list []TfDir
}

func NewTfDirs() *TfDirs {
	return &TfDirs{
		set: make(map[TfDir]bool, 64),
	}
}

func (d *TfDirs) Add(tfDir TfDir) {
	d.set[tfDir] = true
}

func (d *TfDirs) AbsList(baseDir BaseDir) []TfDir {
	result := make([]TfDir, 0, len(d.List()))
	for _, dir := range d.List() {
		result = append(result, filepath.Join(baseDir.Abs(), dir))
	}
	return result
}

func (d *TfDirs) List() []TfDir {
	if d.list != nil {
		return d.list
	}
	return d.generateList()
}

func (d *TfDirs) generateList() []TfDir {
	result := make([]TfDir, 0, len(d.set))
	for dir := range d.set {
		result = append(result, dir)
	}
	sort.Strings(result)
	d.list = result
	return result
}

func (d *TfDirs) String() string {
	return strings.Join(d.List(), " ")
}

func (d *TfDirs) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.List())
}

func (d *TfDirs) ToJson() string {
	return SimpleJsonMarshal(d)
}

type DependentMap struct {
	set map[ModuleDir][]TfDir
}

func NewDependentMap() *DependentMap {
	return &DependentMap{
		set: make(map[ModuleDir][]TfDir, 64),
	}
}

func (m *DependentMap) Add(moduleDir ModuleDir, tfDir TfDir) {
	m.set[moduleDir] = append(m.set[moduleDir], tfDir)
}

func (m *DependentMap) ListTfDirSlice(moduleDir ModuleDir) []TfDir {
	result, _ := m.set[moduleDir]
	return result
}

func (m *DependentMap) IsModule(moduleDir ModuleDir) bool {
	_, ok := m.set[moduleDir]
	return ok
}

func (m *DependentMap) String() string {
	result := ""
	for k, v := range m.set {
		result += fmt.Sprintf("%s:%s\n", k, strings.Join(v, " "))
	}
	return result
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
