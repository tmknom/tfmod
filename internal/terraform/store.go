package terraform

import (
	"log"

	"github.com/tmknom/tfmod/internal/collection"
	"github.com/tmknom/tfmod/internal/dir"
)

type Store interface {
	Save(moduleDir *ModuleDir, stateDir *StateDir)
	List(dirs []*dir.Dir) []string
	Dump()
}

type DependencyStore struct {
	*DependencyGraph
}

func NewDependencyStore() *DependencyStore {
	return &DependencyStore{
		DependencyGraph: NewDependencyGraph(),
	}
}

func (s *DependencyStore) Save(moduleDir *ModuleDir, stateDir *StateDir) {
	s.DependencyGraph.Add(stateDir, moduleDir)
}

func (s *DependencyStore) List(stateDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, stateDir := range stateDirs {
		src := NewStateDir(stateDir.Rel(), stateDir.BaseDir())
		moduleDirs := s.DependencyGraph.ListDst(src)
		for _, moduleDir := range moduleDirs {
			result.Add(moduleDir.Rel())
		}
	}

	return result.Slice()
}

func (s *DependencyStore) Dump() {
	log.Printf("DependencyGraph: %v", s.DependencyGraph)
}

type DependentStore struct {
	*DependentGraph
}

func NewDependentStore() *DependentStore {
	return &DependentStore{
		DependentGraph: NewDependentGraph(),
	}
}

func (s *DependentStore) Save(moduleDir *ModuleDir, stateDir *StateDir) {
	s.DependentGraph.Add(moduleDir, stateDir)
}

func (s *DependentStore) List(moduleDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, moduleDir := range moduleDirs {
		src := NewModuleDir(moduleDir.Rel(), moduleDir.BaseDir())
		stateDirs := s.DependentGraph.ListDst(src)
		for _, stateDir := range stateDirs {
			if !s.DependentGraph.IsModule(stateDir) {
				result.Add(stateDir.Rel())
			}
		}
	}

	return result.Slice()
}

func (s *DependentStore) Dump() {
	log.Printf("DependentGraph: %v", s.DependentGraph)
}
