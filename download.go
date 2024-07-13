package tfmod

import (
	"context"
	"fmt"
	"log"

	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
)

type DownloadRunner struct {
	flags *DownloadFlags
	*IO
}

func NewDownloadRunner(flags *DownloadFlags, io *IO) *DownloadRunner {
	return &DownloadRunner{
		flags: flags,
		IO:    io,
	}
}

type DownloadFlags struct {
	*GlobalFlags
}

func NewDownloadFlags(globalFlags *GlobalFlags) *DownloadFlags {
	return &DownloadFlags{
		GlobalFlags: globalFlags,
	}
}

func (f *DownloadFlags) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

func (r *DownloadRunner) Run(ctx context.Context) error {
	list, err := r.TerraformGet(ctx)
	if err != nil {
		return err
	}
	return format.NewSliceFormatter(r.flags.Format(), list, r.IO.OutWriter).Print()
}

func (r *DownloadRunner) TerraformGet(ctx context.Context) ([]string, error) {
	log.Printf("Runner flags: %#v", r.flags)
	baseDir := r.flags.BaseDir()
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
