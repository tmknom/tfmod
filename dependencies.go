package tfmod

import (
	"fmt"
	"log"
)

type Dependencies struct {
	Store
	BaseDir
	*IO
}

func NewDependencies(baseDir BaseDir, io *IO) *Dependencies {
	return &Dependencies{
		Store:   NewStore(),
		BaseDir: baseDir,
		IO:      io,
	}
}

func (d *Dependencies) Run() error {
	log.Printf("BaseDir: %s", d.BaseDir)

	tfDirs, err := d.BaseDir.ListTfDirs()
	if err != nil {
		return err
	}
	log.Printf("Terraform Directories: %v", tfDirs)

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

	dumped := d.Store.Dump()
	_, err = fmt.Fprintln(d.IO.OutWriter, dumped.String())
	return err
}
