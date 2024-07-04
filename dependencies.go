package tfmod

import (
	"fmt"
	"log"
)

type Dependencies struct {
	Store
	*BaseDir
	*IO
}

func NewDependencies(baseDir *BaseDir, io *IO) *Dependencies {
	return &Dependencies{
		Store:   NewStore(),
		BaseDir: baseDir,
		IO:      io,
	}
}

func (d *Dependencies) Run() error {
	err := NewLoader(d.Store, d.BaseDir, d.IO).Load()
	if err != nil {
		return err
	}
	log.Printf("Load DependentMap from: %v", d.BaseDir)

	result := d.Store.Dump()
	log.Printf("Write stdout from: %#v", result)
	_, err = fmt.Fprintln(d.IO.OutWriter, result.String())
	return err
}
