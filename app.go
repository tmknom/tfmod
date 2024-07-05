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
	*GlobalFlags
}

func NewApp(io *IO, ldflags *Ldflags) *App {
	return &App{
		IO:      io,
		Ldflags: ldflags,
		rootCmd: &cobra.Command{
			Short: "Terraform module mapping tool",
		},
		GlobalFlags: &GlobalFlags{},
	}
}

func (a *App) Run(args []string) error {
	a.prepareCommand(args)

	a.rootCmd.PersistentFlags().StringVarP(&a.GlobalFlags.InputBaseDir, "base-dir", "b", ".", "base directory")

	a.rootCmd.AddCommand(a.newDependenciesCommand())
	a.rootCmd.AddCommand(a.newDependentsCommand())

	return a.rootCmd.Execute()
}

func (a *App) newDependenciesCommand() *cobra.Command {
	flags := NewDependenciesFlags(a.GlobalFlags)
	runner := NewDependencies(flags, NewInMemoryStore(), a.IO)
	command := &cobra.Command{
		Use:   "dependencies",
		Short: "List module dependencies",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run() },
	}
	return command
}

func (a *App) newDependentsCommand() *cobra.Command {
	flags := NewDependentsFlags(a.GlobalFlags)
	runner := NewDependents(flags, NewInMemoryStore(), a.IO)
	command := &cobra.Command{
		Use:   "dependents",
		Short: "List module dependents",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run() },
	}
	command.PersistentFlags().StringSliceVarP(&flags.rawSourceDirs, "sources", "s", []string{}, "source dirs")
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
	cobra.OnInitialize(func() { a.setupLog(a.rootCmd.Name(), args) })

	// setup version option
	version := fmt.Sprintf("%s version %s (%s)", a.Ldflags.Name, a.Ldflags.Version, a.Ldflags.Date)
	a.rootCmd.SetVersionTemplate(version)
}

func (a *App) setupLog(name string, args []string) {
	//log.SetOutput(io.Discard)
	log.SetOutput(os.Stderr)
	log.SetPrefix(fmt.Sprintf("[%s] ", a.Ldflags.Name))
	log.Printf("Start: %s", name)
	log.Printf("Args: %q", args)
	log.Printf("Ldflags: %+v", a.Ldflags)
	log.Printf("Global flags: %#v", a.GlobalFlags)
}

type GlobalFlags struct {
	InputBaseDir string
}

func (f *GlobalFlags) BaseDir() *BaseDir {
	return NewBaseDir(f.InputBaseDir)
}
