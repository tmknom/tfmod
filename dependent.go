package tfmod

import (
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
)

type DependentRunner struct {
	flags *DependentFlags
	terraform.Store
	*IO
}

func NewDependentRunner(flags *DependentFlags, store terraform.Store, io *IO) *DependentRunner {
	return &DependentRunner{
		flags: flags,
		Store: store,
		IO:    io,
	}
}

type DependentFlags struct {
	ModulePaths []string
	*GlobalFlags
}

func NewDependentFlags(globalFlags *GlobalFlags) *DependentFlags {
	return &DependentFlags{
		GlobalFlags: globalFlags,
	}
}

func (f *DependentFlags) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

func (r *DependentRunner) Run() error {
	list, err := r.List()
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(r.flags.Format(), list, r.IO.OutWriter).Print()
}

func (r *DependentRunner) List() ([]string, error) {
	log.Printf("Runner flags: %#v", r.flags)

	baseDir := r.flags.BaseDir()
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

	candidatePaths := r.flags.ModulePaths
	if r.IO.IsPipe() {
		candidatePaths = append(candidatePaths, r.IO.Read()...)
	}
	log.Printf("Candidate paths: %#v", candidatePaths)

	moduleDirs, err := baseDir.FilterDirs(candidatePaths)
	if err != nil {
		return nil, err
	}
	result := r.Store.List(moduleDirs)
	log.Printf("Result: %#v", result)

	return result, nil
}
