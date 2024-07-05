package tfmod

import "log"

type Store interface {
	Save(moduleDir *ModuleDir, tfDir *TfDir)
	ListTfDirs(moduleDirs []string) *Dirs
	ListModuleDirs(stateDirs []string) *Dirs
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

func (s *InMemoryStore) ListTfDirs(moduleDirs []string) *Dirs {
	result := NewDirs()

	for _, moduleDir := range moduleDirs {
		tfDirs := s.DependentMap.ListTfDirSlice(moduleDir)
		for _, tfDir := range tfDirs {
			if !s.DependentMap.IsModule(tfDir.Rel()) {
				result.Add(tfDir.Rel())
			}
		}
	}

	return result
}

func (s *InMemoryStore) ListModuleDirs(stateDirs []string) *Dirs {
	result := NewDirs()

	for _, dir := range stateDirs {
		moduleDirs := s.DependencyMap.ListModuleDirSlice(dir)
		for _, moduleDir := range moduleDirs {
			result.Add(moduleDir.Rel())
		}
	}

	return result
}

func (s *InMemoryStore) Dump() {
	log.Printf("DependencyMap: %v", s.DependencyMap)
	log.Printf("DependentMap: %v", s.DependentMap)
}
