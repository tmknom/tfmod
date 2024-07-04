package tfmod

import (
	"fmt"
	"log"
)

type Dependents struct {
	SliceSourceDirs []string
	Store
	BaseDir
	*IO
}

func NewDependents(io *IO) *Dependents {
	return &Dependents{
		Store: NewStore(),
		IO:    io,
	}
}

func (d *Dependents) InitBaseDir(dir string) {
	if len(d.BaseDir) == 0 {
		d.BaseDir = BaseDir(dir)
	}
}

func (d *Dependents) Run() error {
	err := NewLoader(d.Store, d.BaseDir, d.IO).Load()
	if err != nil {
		return err
	}

	log.Printf("Source Dirs: %v", d.SliceSourceDirs)
	sourceDirs := NewSourceDirs(d.SliceSourceDirs)
	result := d.Store.List(sourceDirs)
	_, err = fmt.Fprintln(d.IO.OutWriter, result.ToJson())

	return nil
}
