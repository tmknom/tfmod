package tfmod

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/tmknom/tfmod/internal/errlib"
)

type Terraform struct{}

func NewTerraform() *Terraform {
	return &Terraform{}
}

func (t *Terraform) ExecuteGetAll(baseDir *BaseDir, dirs *TfDirs, enable bool) error {
	for _, dir := range dirs.AbsList(baseDir) {
		err := t.executeGet(dir, enable)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terraform) executeGet(dir string, enable bool) error {
	if !enable {
		return nil
	}

	cmd := exec.Command("terraform", "get")
	cmd.Dir = dir
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	cmdInfo := fmt.Sprintf("execute: %s (at %s)\n", cmd.String(), cmd.Dir)
	log.Printf(cmdInfo)

	err := cmd.Run()
	if err != nil {
		return errlib.Wrapf(err, "%s\n stdout: %v\n stderr: %v\n", cmdInfo, cmd.Stdout, cmd.Stderr)
	}
	return nil
}

const (
	ModulesPath = ".terraform/modules/modules.json"
)
