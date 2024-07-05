package tfmod

import "log"

type Store interface {
	Save(moduleDir *ModuleDir, tfDir *TfDir)
	List(moduleDirs []string) *TfDirs
	ListModuleDirs(stateDirs []string) *ModuleDirs
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

func (s *InMemoryStore) List(moduleDirs []string) *TfDirs {
	result := NewTfDirs()

	for _, sourceDir := range moduleDirs {
		tfDirs := s.DependentMap.ListTfDirSlice(sourceDir)
		for _, tfDir := range tfDirs {
			if !s.DependentMap.IsModule(tfDir.Rel()) {
				result.Add(tfDir.Rel())
			}
		}
	}

	return result
}

func (s *InMemoryStore) ListModuleDirs(stateDirs []string) *ModuleDirs {
	result := NewModuleDirs()

	for _, dir := range stateDirs {
		moduleDirs := s.DependencyMap.ListModuleDirSlice(dir)
		for _, moduleDir := range moduleDirs {
			//if !s.DependencyMap.IsTf(moduleDir) {
			//	result.Add(moduleDir)
			//}
			result.Add(moduleDir.Rel())
		}
	}

	return result
}

func (s *InMemoryStore) Dump() {
	log.Printf("DependencyMap: %v", s.DependencyMap)
	log.Printf("DependentMap: %v", s.DependentMap)
}
