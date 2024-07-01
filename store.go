package tfmod

type Store interface {
	Save(moduleDir ModuleDir, tfDir TfDir)
	List(sourceDirs SourceDirs) TfDirs
	Dump() DependentMap
}

type InMemoryStore struct {
	DependentMap
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		DependentMap: NewDependentMap(),
	}
}

func (s *InMemoryStore) Save(moduleDir ModuleDir, tfDir TfDir) {
	s.DependentMap[moduleDir] = append(s.DependentMap[moduleDir], tfDir)
}

func (s *InMemoryStore) List(sourceDirs SourceDirs) TfDirs {
	result := NewTfDirs()

	for _, sourceDir := range sourceDirs {
		tfDirs := s.DependentMap[sourceDir]
		for _, tfDir := range tfDirs {
			if !s.DependentMap.IsModule(tfDir) {
				result.Add(tfDir)
			}
		}
	}

	return result
}

func (s *InMemoryStore) Dump() DependentMap {
	return s.DependentMap
}
