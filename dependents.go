package tfmod

import (
	"fmt"
	"log"
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
	rawSourceDirs []string
	*GlobalFlags
}

func NewDependentsFlags(globalFlags *GlobalFlags) *DependentsFlags {
	return &DependentsFlags{
		GlobalFlags: globalFlags,
	}
}

func (d *Dependents) Run() error {
	log.Printf("Runner flags: %#v", d.flags)

	err := NewLoader(d.Store, d.flags.GlobalFlags.BaseDir(), d.flags.GlobalFlags.EnableTf).Load()
	if err != nil {
		return err
	}
	log.Printf("Load DependentMap from: %v", d.flags.BaseDir())

	sourceDirs := NewSourceDirs(d.flags.rawSourceDirs)
	result := d.Store.List(sourceDirs)
	log.Printf("Write stdout from: %#v", result)
	_, err = fmt.Fprintln(d.IO.OutWriter, result.ToJson())
	return err
}
