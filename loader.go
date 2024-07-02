package tfmod

import (
	"log"
)

type Loader struct {
	Store
	BaseDir
	*IO
}

func NewLoader(store Store, baseDir BaseDir, io *IO) *Loader {
	return &Loader{
		Store:   store,
		BaseDir: baseDir,
		IO:      io,
	}
}

func (l *Loader) Load() error {
	log.Printf("BaseDir: %s", l.BaseDir)

	tfDirs, err := l.BaseDir.ListTfDirs()
	if err != nil {
		return err
	}
	log.Printf("Terraform Directories: %v", tfDirs)

	terraform := NewTerraform(l.IO)
	err = terraform.ExecuteGetAll(tfDirs)
	if err != nil {
		return err
	}

	parser := NewParser(l.BaseDir, l.Store)
	err = parser.ParseAll(tfDirs)
	if err != nil {
		return err
	}

	return nil
}
