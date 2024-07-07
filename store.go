package tfmod

import (
	"log"

	"github.com/tmknom/tfmod/internal/collection"
)

type Store interface {
	Save(moduleDir *ModuleDir, tfDir *TfDir)
	ListTfDirs(moduleDirs []string) []string
	ListModuleDirs(stateDirs []string) []string
	Dump()
}

type InMemoryStore struct {
	*DependencyMap
	*DependentMap
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		DependencyMap: NewDependencyMap(),
		DependentMap:  NewDependentMap(),
	}
}

func (s *InMemoryStore) Save(moduleDir *ModuleDir, tfDir *TfDir) {
	s.DependencyMap.Add(tfDir, moduleDir)
	s.DependentMap.Add(moduleDir, tfDir)
}

func (s *InMemoryStore) ListTfDirs(moduleDirs []string) []string {
	result := collection.NewTreeSet()

	for _, moduleDir := range moduleDirs {
		tfDirs := s.DependentMap.ListDst(moduleDir)
		for _, tfDir := range tfDirs {
			if !s.DependentMap.IsModule(tfDir.Rel()) {
				result.Add(tfDir.Rel())
			}
		}
	}

	return result.Slice()
}

func (s *InMemoryStore) ListModuleDirs(stateDirs []string) []string {
	result := collection.NewTreeSet()

	for _, dir := range stateDirs {
		moduleDirs := s.DependencyMap.ListDst(dir)
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
