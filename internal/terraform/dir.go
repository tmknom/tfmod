package terraform

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

type DependentGraph struct {
	*dir.Graph[ModuleDir, TfDir]
}

func NewDependentGraph() *DependentGraph {
	return &DependentGraph{
		Graph: dir.NewGraph[ModuleDir, TfDir](),
	}
}

func (m *DependentGraph) IsModule(dir string) bool {
	return m.Include(dir)
}

type DependencyGraph struct {
	*dir.Graph[TfDir, ModuleDir]
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Graph: dir.NewGraph[TfDir, ModuleDir](),
	}
}
