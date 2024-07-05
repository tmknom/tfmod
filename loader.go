package tfmod

import (
	"log"
)

type Loader struct {
	Store
	*BaseDir
	*IO
}

func NewLoader(store Store, baseDir *BaseDir, io *IO) *Loader {
	return &Loader{
		Store:   store,
		BaseDir: baseDir,
		IO:      io,
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
	err = terraform.ExecuteGetAll(*l.BaseDir, tfDirs)
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
