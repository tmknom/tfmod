package tfmod

import (
	"fmt"
	"log"
)

type Dependents struct {
	SliceSourceDirs []string
	Store
	*BaseDir
	*IO
}

func NewDependents(io *IO) *Dependents {
	return &Dependents{
		Store: NewStore(),
		IO:    io,
	}
}

func (d *Dependents) InitBaseDir(dir string) {
	if d.BaseDir == nil {
		d.BaseDir = NewBaseDir(dir)
	}
}

func (d *Dependents) Run() error {
	log.Printf("Source Dirs: %v", d.SliceSourceDirs)

	err := NewLoader(d.Store, d.BaseDir, d.IO).Load()
	if err != nil {
		return err
	}
	log.Printf("Load DependentMap from: %v", d.BaseDir)

	sourceDirs := NewSourceDirs(d.SliceSourceDirs)
	result := d.Store.List(sourceDirs)
	log.Printf("Write stdout from: %#v", result)
	_, err = fmt.Fprintln(d.IO.OutWriter, result.ToJson())
	return err
}
