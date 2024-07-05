package tfmod

import (
	"log"
)

type Loader struct {
	Store
	*BaseDir
	enableTf bool
}

func NewLoader(store Store, baseDir *BaseDir, enableTf bool) *Loader {
	return &Loader{
		Store:    store,
		BaseDir:  baseDir,
		enableTf: enableTf,
	}
}

func (l *Loader) Load() error {
	log.Printf("BaseDir: %v", l.BaseDir)

	sourceDirs, err := l.BaseDir.GenerateSourceDirs()
	if err != nil {
		return err
	}
	log.Printf("Source dirs: %v", sourceDirs)

	terraform := NewTerraform()
	err = terraform.ExecuteGetAll(sourceDirs, l.enableTf)
	if err != nil {
		return err
	}
	log.Printf("Execute terraform get to: %v", sourceDirs)

	parser := NewParser(l.BaseDir, l.Store)
	err = parser.ParseAll(sourceDirs)
	if err != nil {
		return err
	}
	log.Printf("Parse to: %v", sourceDirs)

	return nil
}
