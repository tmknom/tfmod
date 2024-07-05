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
			log.Fatalf("%+v", errlib.Wrapf(err, "invalid current dir: %s", dir))
		}
		dir = filepath.Join(currentDir, dir)
	}

	result, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "invalid base dir: %s", dir))
	}
	return result
}

func (d *BaseDir) GenerateSourceDirs() ([]*SourceDir, error) {
	sourceDirs := make([]*SourceDir, 0, 64)

	err := filepath.WalkDir(d.Abs(), func(absFilepath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return errlib.Wrapf(err, "invalid base dir: %#v", d)
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".tf" {
			target := NewSourceDir(filepath.Dir(absFilepath), d)
			sourceDirs = append(sourceDirs, target)
		}
		return nil
	})
	return sourceDirs, err
}

type SourceDir struct {
	abs     string
	baseDir *BaseDir
}

func NewSourceDir(abs string, baseDir *BaseDir) *SourceDir {
	return &SourceDir{
		abs:     abs,
		baseDir: baseDir,
	}
}

func (d *SourceDir) Abs() string {
	return d.abs
}

func (d *SourceDir) Rel() string {
	rel, _ := filepath.Rel(d.baseDir.Abs(), d.Abs())
	return rel
}

func (d *SourceDir) ToTfDir() *TfDir {
	return NewTfDir(d.Rel(), d.baseDir)
}

func (d *SourceDir) AbsBaseDir() string {
	return d.baseDir.Abs()
}

func (d *SourceDir) String() string {
	return d.Rel()
}

type ModuleDir struct {
	raw     string
	baseDir *BaseDir
}

func NewModuleDir(raw string, baseDir *BaseDir) *ModuleDir {
	return &ModuleDir{
		raw:     raw,
		baseDir: baseDir,
	}
}

func (d *ModuleDir) Rel() string {
	return d.raw
}

func (d *ModuleDir) String() string {
	return d.Rel()
}

type ModuleDirs struct {
	set  map[string]bool
	list []string
}

type TfDir struct {
	raw     string
	baseDir *BaseDir
}

func NewTfDir(raw string, baseDir *BaseDir) *TfDir {
	return &TfDir{
		raw:     raw,
		baseDir: baseDir,
	}
}

func (d *TfDir) Rel() string {
	return d.raw
}

func (d *TfDir) String() string {
	return d.Rel()
}

func NewModuleDirs() *ModuleDirs {
	return &ModuleDirs{
		set: make(map[string]bool, 64),
	}
}

func (d *ModuleDirs) Add(dir string) {
	d.set[dir] = true
}

func (d *ModuleDirs) List() []string {
	if d.list != nil {
		return d.list
	}
	return d.generateList()
}

func (d *ModuleDirs) generateList() []string {
	result := make([]string, 0, len(d.set))
	for dir := range d.set {
		result = append(result, dir)
	}
	sort.Strings(result)
	d.list = result
	return result
}

func (d *ModuleDirs) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.List())
}

func (d *ModuleDirs) ToJson() string {
	return SimpleJsonMarshal(d)
}

type TfDirs struct {
	set  map[string]bool
	list []string
}

func NewTfDirs() *TfDirs {
	return &TfDirs{
		set: make(map[string]bool, 64),
	}
}

func (d *TfDirs) Add(tfDir string) {
	d.set[tfDir] = true
}

func (d *TfDirs) List() []string {
	if d.list != nil {
		return d.list
	}
	return d.generateList()
}

func (d *TfDirs) generateList() []string {
	result := make([]string, 0, len(d.set))
	for dir := range d.set {
		result = append(result, dir)
	}
	sort.Strings(result)
	d.list = result
	return result
}

func (d *TfDirs) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.List())
}

func (d *TfDirs) ToJson() string {
	return SimpleJsonMarshal(d)
}

type DependentMap struct {
	set map[string][]*TfDir
}

func NewDependentMap() *DependentMap {
	return &DependentMap{
		set: make(map[string][]*TfDir, 64),
	}
}

func (m *DependentMap) Add(moduleDir *ModuleDir, tfDir *TfDir) {
	key := moduleDir.Rel()
	m.set[key] = append(m.set[key], tfDir)
}

func (m *DependentMap) ListTfDirSlice(moduleDir string) []*TfDir {
	result, _ := m.set[moduleDir]
	return result
}

func (m *DependentMap) IsModule(moduleDir string) bool {
	_, ok := m.set[moduleDir]
	return ok
}

func (m *DependentMap) String() string {
	items := make([]string, 0, len(m.set))
	for key, tfDirs := range m.set {
		strDirs := make([]string, 0, len(tfDirs))
		for _, tfDir := range tfDirs {
			strDirs = append(strDirs, tfDir.Rel())
		}
		items = append(items, fmt.Sprintf("%s:%v", key, strDirs))
	}
	return strings.Join(items, ", ")
}

func (m *DependentMap) ToJson() string {
	return SimpleJsonMarshal(m)
}

type DependencyMap struct {
	set map[string][]*ModuleDir
}

func NewDependencyMap() *DependencyMap {
	return &DependencyMap{
		set: make(map[string][]*ModuleDir, 64),
	}
}

func (m *DependencyMap) Add(tfDir *TfDir, moduleDir *ModuleDir) {
	key := tfDir.Rel()
	m.set[key] = append(m.set[key], moduleDir)
}

func (m *DependencyMap) ListModuleDirSlice(tfDir string) []*ModuleDir {
	result, _ := m.set[tfDir]
	return result
}

func (m *DependencyMap) String() string {
	items := make([]string, 0, len(m.set))
	for key, moduleDirs := range m.set {
		strDirs := make([]string, 0, len(moduleDirs))
		for _, moduleDir := range moduleDirs {
			strDirs = append(strDirs, moduleDir.Rel())
		}
		items = append(items, fmt.Sprintf("%s:%v", key, strDirs))
	}
	return strings.Join(items, ", ")
}

func (m *DependencyMap) ToJson() string {
	return SimpleJsonMarshal(m)
}

func SimpleJsonMarshal(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "invalid json: %#v", v))
	}
	return string(bytes)
}
