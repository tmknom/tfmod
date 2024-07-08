package tfmod

import (
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
)

type DependenciesRunner struct {
	flags *DependenciesFlags
	Store
	*IO
}

func NewDependenciesRunner(flags *DependenciesFlags, store Store, io *IO) *DependenciesRunner {
	return &DependenciesRunner{
		flags: flags,
		Store: store,
		IO:    io,
	}
}

type DependenciesFlags struct {
	StateDirs []string
	*GlobalFlags
}

func NewDependenciesFlags(globalFlags *GlobalFlags) *DependenciesFlags {
	return &DependenciesFlags{
		GlobalFlags: globalFlags,
	}
}

func (f *DependenciesFlags) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

func (r *DependenciesRunner) Run() error {
	list, err := r.List()
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(r.flags.Format, list, r.IO.OutWriter).Print()
}

func (r *DependenciesRunner) List() ([]string, error) {
	log.Printf("Runner flags: %#v", r.flags)

	baseDir := r.flags.GetBaseDir()
	filter := terraform.NewFilter(baseDir)
	sourceDirs, err := filter.SubDirs()
	if err != nil {
		return nil, err
	}

	parser := NewParser(r.Store)
	err = parser.ParseAll(sourceDirs)
	if err != nil {
		return nil, err
	}

	r.Store.Dump()
	result := r.Store.ListModuleDirs(baseDir.ConvertDirs(r.flags.StateDirs))
	log.Printf("Result: %#v", result)

	return result, nil
}
