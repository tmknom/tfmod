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

func (t *Terraform) ExecuteGetAll(dirs *TfDirs) error {
	for _, dir := range dirs.List() {
		err := t.executeGet(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terraform) executeGet(dir string) error {
	if SuppressTerraform {
		return nil
	}

	cmd := exec.Command("terraform", "get")
	cmd.Dir = dir
	cmd.Stdout = t.OutWriter
	cmd.Stderr = t.ErrWriter

	log.Printf("execute: %s (at %s)\n", cmd.String(), cmd.Dir)
	return cmd.Run()
}

const (
	ModulesPath       = ".terraform/modules/modules.json"
	SuppressTerraform = true
)
