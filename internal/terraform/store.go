package terraform

import (
	"log"

	"github.com/tmknom/tfmod/internal/collection"
	"github.com/tmknom/tfmod/internal/dir"
)

type Store interface {
	Save(moduleDir *ModuleDir, tfDir *TfDir)
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

func (s *DependencyStore) Save(moduleDir *ModuleDir, tfDir *TfDir) {
	s.DependencyGraph.Add(tfDir, moduleDir)
}

func (s *DependencyStore) List(stateDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, stateDir := range stateDirs {
		src := NewTfDir(stateDir.Rel(), stateDir.BaseDir())
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

func (s *DependentStore) Save(moduleDir *ModuleDir, tfDir *TfDir) {
	s.DependentGraph.Add(moduleDir, tfDir)
}

func (s *DependentStore) List(moduleDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, moduleDir := range moduleDirs {
		src := NewModuleDir(moduleDir.Rel(), moduleDir.BaseDir())
		tfDirs := s.DependentGraph.ListDst(src)
		for _, tfDir := range tfDirs {
			if !s.DependentGraph.IsModule(tfDir) {
				result.Add(tfDir.Rel())
			}
		}
	}

	return result.Slice()
}

func (s *DependentStore) Dump() {
	log.Printf("DependentGraph: %v", s.DependentGraph)
}
