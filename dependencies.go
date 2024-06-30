package tfmod

import (
	"encoding/json"
	"fmt"
	"log"
)

type BaseDir = string
type TfDirs map[string]bool

type Dependencies struct {
	BaseDir
	*IO
}

func NewDependencies(baseDir BaseDir, io *IO) *Dependencies {
	return &Dependencies{
		BaseDir: baseDir,
		IO:      io,
	}
}

func (d *Dependencies) Run() error {
	log.Printf("BaseDir: %s", d.BaseDir)

	tfDirs, err := NewFilter(d.BaseDir).FilterTfDirs()
	if err != nil {
		return err
	}
	log.Printf("Terraform Directories: %v", tfDirs)

	terraform := NewTerraform(d.IO)
	err = terraform.RunGetAll(tfDirs)
	if err != nil {
		return err
	}

	parser := NewParser(d.BaseDir)
	err = parser.ParseAll(tfDirs)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(parser.Mapping)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}
