package tfmod

import (
	"context"
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
)

type GetRunner struct {
	flags *GetFlags
	*IO
}

func NewGetRunner(flags *GetFlags, io *IO) *GetRunner {
	return &GetRunner{
		flags: flags,
		IO:    io,
	}
}

type GetFlags struct {
	*GlobalFlags
}

func NewGetFlags(globalFlags *GlobalFlags) *GetFlags {
	return &GetFlags{
		GlobalFlags: globalFlags,
	}
}

func (f *GetFlags) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

func (r *GetRunner) Run(ctx context.Context) error {
	list, err := r.TerraformGet(ctx)
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(r.flags.Format(), list, r.IO.OutWriter).Print()
}

func (r *GetRunner) TerraformGet(ctx context.Context) ([]string, error) {
	log.Printf("Runner flags: %#v", r.flags)
	baseDir := r.flags.GetBaseDir()
	filter := terraform.NewFilter(baseDir)
	sourceDirs, err := filter.SubDirs()
	if err != nil {
		return nil, err
	}

	terraformCommand := terraform.NewCommand()
	err = terraformCommand.GetAll(ctx, sourceDirs)
	if err != nil {
		return nil, err
	}

	log.Printf("Result: %#v", sourceDirs)
	result := make([]string, 0, len(sourceDirs))
	for _, srcDir := range sourceDirs {
		result = append(result, srcDir.Rel())
	}
	return result, nil
}
