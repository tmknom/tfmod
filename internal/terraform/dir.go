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

type StateDir struct {
	*dir.Dir
}

func NewStateDir(raw string, baseDir *dir.BaseDir) *StateDir {
	return &StateDir{
		Dir: baseDir.CreateDir(raw),
	}
}

type DependentGraph struct {
	*dir.Graph[ModuleDir, StateDir]
}

func NewDependentGraph() *DependentGraph {
	return &DependentGraph{
		Graph: dir.NewGraph[ModuleDir, StateDir](),
	}
}

func (m *DependentGraph) IsModule(src dir.Path) bool {
	return m.Include(src)
}

type DependencyGraph struct {
	*dir.Graph[StateDir, ModuleDir]
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Graph: dir.NewGraph[StateDir, ModuleDir](),
	}
}
