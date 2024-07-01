package tfmod

import (
	"fmt"
	"log"
)

type Dependents struct {
	SliceSourceDirs []string
	Store
	BaseDir
	*IO
}

func NewDependents(baseDir BaseDir, io *IO) *Dependents {
	return &Dependents{
		Store:   NewStore(),
		BaseDir: baseDir,
		IO:      io,
	}
}

func (d *Dependents) Run() error {
	log.Printf("Source Dirs: %v", d.SliceSourceDirs)
	log.Printf("BaseDir: %s", d.BaseDir)

	tfDirs, err := d.BaseDir.ListTfDirs()
	if err != nil {
		return err
	}
	log.Printf("Terraform Dirs: %s", tfDirs.String())

	terraform := NewTerraform(d.IO)
	err = terraform.ExecuteGetAll(tfDirs)
	if err != nil {
		return err
	}

	parser := NewParser(d.BaseDir, d.Store)
	err = parser.ParseAll(tfDirs)
	if err != nil {
		return err
	}

	sourceDirs := NewSourceDirs(d.SliceSourceDirs)
	result := d.Store.List(sourceDirs)
	_, err = fmt.Fprintln(d.IO.OutWriter, result.ToJson())

	return nil
}
