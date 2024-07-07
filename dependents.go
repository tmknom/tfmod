package tfmod

import (
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
)

type Dependents struct {
	flags *DependentsFlags
	Store
	*IO
}

func NewDependents(flags *DependentsFlags, store Store, io *IO) *Dependents {
	return &Dependents{
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

func (d *Dependents) Run() error {
	list, err := d.List()
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(d.flags.Format, list, d.IO.OutWriter).Print()
}

func (d *Dependents) List() ([]string, error) {
	log.Printf("Runner flags: %#v", d.flags)

	err := NewLoader(d.Store, d.flags.GetBaseDir(), d.flags.EnableTf).Load()
	if err != nil {
		return nil, err
	}

	result := d.Store.ListTfDirs(d.flags.ModuleDirs)
	log.Printf("Write stdout from: %#v", result)

	return result, nil
}
