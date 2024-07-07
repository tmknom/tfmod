package tfmod

import (
	"github.com/tmknom/tfmod/internal/dir"
)

type ModuleDir struct {
	*dir.Dir
}

func NewModuleDir(raw string, baseDir *dir.BaseDir) *ModuleDir {
	return &ModuleDir{
		Dir: baseDir.CreateDir(raw),
	}
}

type TfDir struct {
	*dir.Dir
}

func NewTfDir(raw string, baseDir *dir.BaseDir) *TfDir {
	return &TfDir{
		Dir: baseDir.CreateDir(raw),
	}
}

type DependentMap struct {
	*dir.Graph[ModuleDir, TfDir]
}

func NewDependentMap() *DependentMap {
	return &DependentMap{
		Graph: dir.NewGraph[ModuleDir, TfDir](),
	}
}

func (m *DependentMap) IsModule(dir string) bool {
	return m.Include(dir)
}

type DependencyMap struct {
	*dir.Graph[TfDir, ModuleDir]
}

func NewDependencyMap() *DependencyMap {
	return &DependencyMap{
		Graph: dir.NewGraph[TfDir, ModuleDir](),
	}
}
