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
	rootCmd *cobra.Command
}

func NewApp(io *IO, ldflags *Ldflags) *App {
	return &App{
		IO:      io,
		Ldflags: ldflags,
		rootCmd: &cobra.Command{
			Short: "Terraform module mapping tool",
		},
	}
}

func (a *App) Run(args []string) error {
	a.prepareCommand(args)

	a.rootCmd.AddCommand(a.newDependenciesCommand())
	a.rootCmd.AddCommand(a.newDependentsCommand())

	return a.rootCmd.Execute()
}

func (a *App) newDependenciesCommand() *cobra.Command {
	currentDir, _ := os.Getwd()
	runner := NewDependencies(BaseDir(currentDir), a.IO)
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

func (a *App) newDependentsCommand() *cobra.Command {
	currentDir, _ := os.Getwd()
	runner := NewDependents(BaseDir(currentDir), a.IO)
	command := &cobra.Command{
		Use:   "dependents",
		Short: "List module dependents",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("start: %s", cmd.Name())
			return runner.Run()
		},
	}
	command.PersistentFlags().StringSliceVarP(&runner.SliceSourceDirs, "sources", "s", []string{}, "source dirs")
	return command
}

func (a *App) prepareCommand(args []string) {
	// set ldflags
	a.rootCmd.Use = a.Ldflags.Name
	a.rootCmd.Version = a.Ldflags.Version

	// override default settings
	a.rootCmd.SetArgs(args)
	a.rootCmd.SetIn(a.IO.InReader)
	a.rootCmd.SetOut(a.IO.OutWriter)
	a.rootCmd.SetErr(a.IO.ErrWriter)

	// setup log
	cobra.OnInitialize(func() { a.setupLog(args) })

	// setup version option
	version := fmt.Sprintf("%s version %s (%s)", a.Ldflags.Name, a.Ldflags.Version, a.Ldflags.Date)
	a.rootCmd.SetVersionTemplate(version)
}

func (a *App) setupLog(args []string) {
	//log.SetOutput(io.Discard)
	log.SetOutput(os.Stderr)
	log.SetPrefix(fmt.Sprintf("[%s] ", a.Ldflags.Name))
	log.Printf("args: %q", args)
	log.Printf("ldflags: %+v", a.Ldflags)
}
