package tfmod

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type IO struct {
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
}

type Ldflags struct {
	Name    string
	Version string
	Commit  string
	Date    string
}

type App struct {
	*IO
	*Ldflags
}

func NewApp(io *IO, ldflags *Ldflags) *App {
	return &App{
		IO:      io,
		Ldflags: ldflags,
	}
}

func (a *App) Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:     a.Ldflags.Name,
		Short:   "Terraform module mapping tool",
		Version: a.Ldflags.Version,
	}

	// override default settings
	rootCmd.SetArgs(args)
	rootCmd.SetIn(a.IO.InReader)
	rootCmd.SetOut(a.IO.OutWriter)
	rootCmd.SetErr(a.IO.ErrWriter)

	// setup log
	cobra.OnInitialize(func() { a.setupLog(args) })

	// setup version option
	version := fmt.Sprintf("%s version %s (%s)", a.Ldflags.Name, a.Ldflags.Version, a.Ldflags.Date)
	rootCmd.SetVersionTemplate(version)

	// setup commands
	rootCmd.AddCommand(a.newGenerateCommand())

	return rootCmd.Execute()
}

func (a *App) newGenerateCommand() *cobra.Command {
	currentDir, _ := os.Getwd()
	runner := NewDependencies(currentDir, a.IO)
	command := &cobra.Command{
		Use:   "dependencies",
		Short: "List module dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("start: %s", cmd.Name())
			return runner.Run()
		},
	}
	return command
}

func (a *App) setupLog(args []string) {
	//log.SetOutput(io.Discard)
	log.SetOutput(os.Stderr)
	log.SetPrefix(fmt.Sprintf("[%s] ", a.Ldflags.Name))
	log.Printf("args: %q", args)
	log.Printf("ldflags: %+v", a.Ldflags)
}
