package tfmod

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/errlib"
)

type Terraform struct {
	baseDir *dir.BaseDir
	enable  bool
}

func NewTerraform(baseDir *dir.BaseDir, enable bool) *Terraform {
	return &Terraform{
		baseDir: baseDir,
		enable:  enable,
	}
}

func (t *Terraform) GetAll() ([]*dir.Dir, error) {
	sourceDirs, err := t.baseDir.FilterSubDirs(".tf", filepath.Dir(TerraformModulesPath))
	if err != nil {
		return nil, err
	}
	log.Printf("Source dirs: %v", sourceDirs)

	err = t.executeGetAll(sourceDirs)
	if err != nil {
		return nil, err
	}
	return sourceDirs, nil
}

func (t *Terraform) executeGetAll(workDirs []*dir.Dir) error {
	for _, workDir := range workDirs {
		err := t.executeGet(workDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terraform) executeGet(workDir *dir.Dir) error {
	cmd := exec.Command("terraform", "get")
	cmd.Dir = workDir.Abs()
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	info := fmt.Sprintf("%s (at %s)\n", cmd.String(), cmd.Dir)
	if !t.enable {
		log.Printf(fmt.Sprintf("skip: %s", info))
		return nil
	}
	log.Printf(fmt.Sprintf("execute: %s", info))

	err := cmd.Run()
	if err != nil {
		return errlib.Wrapf(err, "%s\n stdout: %v\n stderr: %v\n", info, cmd.Stdout, cmd.Stderr)
	}
	return nil
}
