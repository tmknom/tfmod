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
	*terraform.DependencyMap
	*terraform.DependentMap
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		DependencyMap: terraform.NewDependencyMap(),
		DependentMap:  terraform.NewDependentMap(),
	}
}

func (s *InMemoryStore) Save(moduleDir *terraform.ModuleDir, tfDir *terraform.TfDir) {
	s.DependencyMap.Add(tfDir, moduleDir)
	s.DependentMap.Add(moduleDir, tfDir)
}

func (s *InMemoryStore) ListTfDirs(moduleDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, moduleDir := range moduleDirs {
		tfDirs := s.DependentMap.ListDst(moduleDir.Rel())
		for _, tfDir := range tfDirs {
			if !s.DependentMap.IsModule(tfDir.Rel()) {
				result.Add(tfDir.Rel())
			}
		}
	}

	return result.Slice()
}

func (s *InMemoryStore) ListModuleDirs(stateDirs []*dir.Dir) []string {
	result := collection.NewTreeSet()

	for _, stateDir := range stateDirs {
		moduleDirs := s.DependencyMap.ListDst(stateDir.Rel())
		for _, moduleDir := range moduleDirs {
			result.Add(moduleDir.Rel())
		}
	}

	return result.Slice()
}

func (s *InMemoryStore) Dump() {
	log.Printf("DependencyMap: %v", s.DependencyMap)
	log.Printf("DependentMap: %v", s.DependentMap)
}
