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

	tfDirs, err := l.BaseDir.ListTfDirs()
	if err != nil {
		return err
	}
	log.Printf("Generate tfdirs from: %v", l.BaseDir)

	terraform := NewTerraform()
	err = terraform.ExecuteGetAll(l.BaseDir, tfDirs, l.enableTf)
	if err != nil {
		return err
	}
	log.Printf("Execute terraform get to: %v", tfDirs)

	parser := NewParser(l.BaseDir, l.Store)
	err = parser.ParseAll(tfDirs)
	if err != nil {
		return err
	}
	log.Printf("Parse to: %v", tfDirs)

	return nil
}
