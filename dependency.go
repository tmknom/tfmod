package tfmod

import (
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
)

type DependencyRunner struct {
	flags *DependencyFlags
	terraform.Store
	*IO
}

func NewDependencyRunner(flags *DependencyFlags, store terraform.Store, io *IO) *DependencyRunner {
	return &DependencyRunner{
		flags: flags,
		Store: store,
		IO:    io,
	}
}

type DependencyFlags struct {
	StatePaths []string
	*GlobalFlags
}

func NewDependencyFlags(globalFlags *GlobalFlags) *DependencyFlags {
	return &DependencyFlags{
		GlobalFlags: globalFlags,
	}
}

func (f *DependencyFlags) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

func (r *DependencyRunner) Run() error {
	list, err := r.List()
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(r.flags.Format(), list, r.IO.OutWriter).Print()
}

func (r *DependencyRunner) List() ([]string, error) {
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

	candidatePaths := r.flags.StatePaths
	if r.IO.IsPipe() {
		candidatePaths = append(candidatePaths, r.IO.Read()...)
	}
	log.Printf("Candidate paths: %#v", candidatePaths)

	stateDirs, err := baseDir.FilterDirs(candidatePaths)
	if err != nil {
		return nil, err
	}
	result := r.Store.List(stateDirs)
	log.Printf("Result: %#v", result)

	return result, nil
}
