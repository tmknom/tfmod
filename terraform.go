package tfmod

import (
	"log"
	"os/exec"
)

type Terraform struct {
	*IO
}

func NewTerraform(io *IO) *Terraform {
	return &Terraform{
		IO: io,
	}
}

func (t *Terraform) RunGetAll(dirs TfDirs) error {
	for dir := range dirs {
		err := t.runGet(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terraform) runGet(dir string) error {
	if SuppressTerraform {
		return nil
	}

	log.Printf("terraform get: %s\n", dir)
	cmd := exec.Command("terraform", "get")
	cmd.Dir = dir
	cmd.Stdout = t.OutWriter
	cmd.Stderr = t.ErrWriter
	return cmd.Run()
}

const (
	ModulesPath       = ".terraform/modules/modules.json"
	SuppressTerraform = true
)
