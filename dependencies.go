package tfmod

import (
	"fmt"
	"log"
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
	*GlobalFlags
}

func NewDependenciesFlags(globalFlags *GlobalFlags) *DependenciesFlags {
	return &DependenciesFlags{
		GlobalFlags: globalFlags,
	}
}

func (d *Dependencies) Run() error {
	log.Printf("Runner flags: %#v", d.flags)

	err := NewLoader(d.Store, d.flags.GlobalFlags.BaseDir(), d.flags.GlobalFlags.EnableTf).Load()
	if err != nil {
		return err
	}
	log.Printf("Load DependentMap from: %v", d.flags.BaseDir())

	result := d.Store.Dump()
	log.Printf("Write stdout from: %#v", result)
	_, err = fmt.Fprintln(d.IO.OutWriter, result.String())
	return err
}
