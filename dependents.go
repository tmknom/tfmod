package tfmod

import (
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
)

type DependentsRunner struct {
	flags *DependentsFlags
	terraform.Store
	*IO
}

func NewDependentsRunner(flags *DependentsFlags, store terraform.Store, io *IO) *DependentsRunner {
	return &DependentsRunner{
		flags: flags,
		Store: store,
		IO:    io,
	}
}

type DependentsFlags struct {
	ModuleDirs []string
	*GlobalFlags
}

func NewDependentsFlags(globalFlags *GlobalFlags) *DependentsFlags {
	return &DependentsFlags{
		GlobalFlags: globalFlags,
	}
}

func (f *DependentsFlags) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

func (r *DependentsRunner) Run() error {
	list, err := r.List()
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(r.flags.Format(), list, r.IO.OutWriter).Print()
}

func (r *DependentsRunner) List() ([]string, error) {
	log.Printf("Runner flags: %#v", r.flags)

	baseDir := r.flags.GetBaseDir()
	filter := terraform.NewFilter(baseDir)
	sourceDirs, err := filter.SubDirs()
	if err != nil {
		return nil, err
	}

	parser := terraform.NewParser(r.Store)
	err = parser.ParseAll(sourceDirs)
	if err != nil {
		return nil, err
	}

	r.Store.Dump()

	moduleDirs, err := baseDir.ConvertDirs(r.flags.ModuleDirs)
	if err != nil {
		return nil, err
	}
	result := r.Store.List(moduleDirs)
	log.Printf("Result: %#v", result)

	return result, nil
}
