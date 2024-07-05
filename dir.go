package tfmod

import (
	"fmt"
	"strings"

	"github.com/tmknom/tfmod/internal/collection"
	"github.com/tmknom/tfmod/internal/dir"
)

type Dirs struct {
	*collection.TreeSet
}

func NewDirs() *Dirs {
	return &Dirs{
		TreeSet: collection.NewTreeSet(),
	}
}

type SourceDir struct {
	*dir.Dir
}

func NewSourceDir(raw string, baseDir *dir.BaseDir) *SourceDir {
	return &SourceDir{
		Dir: dir.NewDir(raw, baseDir),
	}
}

func (d *SourceDir) ToTfDir() *TfDir {
	return NewTfDir(d.Rel(), d.BaseDir())
}

type ModuleDir struct {
	*dir.Dir
}

func NewModuleDir(raw string, baseDir *dir.BaseDir) *ModuleDir {
	return &ModuleDir{
		Dir: dir.NewDir(raw, baseDir),
	}
}

type TfDir struct {
	*dir.Dir
}

func NewTfDir(raw string, baseDir *dir.BaseDir) *TfDir {
	return &TfDir{
		Dir: dir.NewDir(raw, baseDir),
	}
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
