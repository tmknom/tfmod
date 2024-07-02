package tfmod

import (
	"fmt"
)

type Dependencies struct {
	Store
	BaseDir
	*IO
}

func NewDependencies(baseDir BaseDir, io *IO) *Dependencies {
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

	dumped := d.Store.Dump()
	_, err = fmt.Fprintln(d.IO.OutWriter, dumped.String())
	return err
}
