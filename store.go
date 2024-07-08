package tfmod

import (
	"log"

	"github.com/tmknom/tfmod/internal/collection"
	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/terraform"
)

type Store interface {
	Save(moduleDir *terraform.ModuleDir, tfDir *terraform.TfDir)
	ListTfDirs(moduleDirs []*dir.Dir) []string
	ListModuleDirs(stateDirs []*dir.Dir) []string
	Dump()
}

type InMemoryStore struct {
	*terraform.DependencyGraph
	*terraform.DependentGraph
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		DependencyGraph: terraform.NewDependencyGraph(),
		DependentGraph:  terraform.NewDependentGraph(),
	}
}

func (s *InMemoryStore) Save(moduleDir *terraform.ModuleDir, tfDir *terraform.TfDir) {
	s.DependencyGraph.Add(tfDir, moduleDir)
	s.DependentGraph.Add(moduleDir, tfDir)
}

func (s *InMemoryStore) ListTfDirs(moduleDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, moduleDir := range moduleDirs {
		tfDirs := s.DependentGraph.ListDst(moduleDir.Rel())
		for _, tfDir := range tfDirs {
			if !s.DependentGraph.IsModule(tfDir.Rel()) {
				result.Add(tfDir.Rel())
			}
		}
	}

	return result.Slice()
}

func (s *InMemoryStore) ListModuleDirs(stateDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, stateDir := range stateDirs {
		moduleDirs := s.DependencyGraph.ListDst(stateDir.Rel())
		for _, moduleDir := range moduleDirs {
			result.Add(moduleDir.Rel())
		}
	}

	return result.Slice()
}

func (s *InMemoryStore) Dump() {
	log.Printf("DependencyGraph: %v", s.DependencyGraph)
	log.Printf("DependentGraph: %v", s.DependentGraph)
}
