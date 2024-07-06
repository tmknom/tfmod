package tfmod

import (
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
)

type Dependencies struct {
	flags *DependenciesFlags
	Store
	*IO
}

func NewDependencies(flags *DependenciesFlags, store Store, io *IO) *Dependencies {
	return &Dependencies{
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

func (d *Dependencies) Run() error {
	list, err := d.List()
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(d.flags.Format, list, d.IO.OutWriter).Print()
}

func (d *Dependencies) List() ([]string, error) {
	log.Printf("Runner flags: %#v", d.flags)

	err := NewLoader(d.Store, d.flags.GetBaseDir(), d.flags.EnableTf).Load()
	if err != nil {
		return nil, err
	}

	result := d.Store.ListModuleDirs(d.flags.StateDirs)
	log.Printf("Result: %#v", result)

	return result, nil
}
