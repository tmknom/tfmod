package tfmod

import (
	"fmt"

	"github.com/tmknom/tfmod/internal/dir"
)

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
	*DirMap[ModuleDir, TfDir]
}

func NewDependentMap() *DependentMap {
	return &DependentMap{
		DirMap: NewDirMap[ModuleDir, TfDir](),
	}
}

func (m *DependentMap) IsModule(dir string) bool {
	return m.Include(dir)
}

type DependencyMap struct {
	*DirMap[TfDir, ModuleDir]
}

func NewDependencyMap() *DependencyMap {
	return &DependencyMap{
		DirMap: NewDirMap[TfDir, ModuleDir](),
	}
}

type DirMap[S, D dir.Path] struct {
	items map[string][]*D
}

func NewDirMap[S, D dir.Path]() *DirMap[S, D] {
	return &DirMap[S, D]{
		items: make(map[string][]*D, 64),
	}
}

func (m *DirMap[S, D]) Add(src *S, dst *D) {
	key := (*src).Rel()
	m.items[key] = append(m.items[key], dst)
}

func (m *DirMap[S, D]) Include(src string) bool {
	_, ok := m.items[src]
	return ok
}

func (m *DirMap[S, D]) ListDst(src string) []*D {
	result, _ := m.items[src]
	return result
}

func (m *DirMap[S, D]) String() string {
	return fmt.Sprintf("%v", m.items)
}
