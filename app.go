package tfmod

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/format"
	"github.com/tmknom/tfmod/internal/terraform"
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
	a.rootCmd.SetContext(context.Background())

	a.rootCmd.PersistentFlags().StringVarP(&a.GlobalFlags.base, "base", "b", ".", "The base directory that contains tf files")
	a.rootCmd.PersistentFlags().StringVarP(&a.GlobalFlags.format, "format", "f", format.TextFormat, fmt.Sprintf("Format output by: {%s}", format.SupportType()))
	a.rootCmd.PersistentFlags().BoolVar(&a.GlobalFlags.debug, "debug", false, "Show debugging output")

	a.rootCmd.AddCommand(a.newGetCommand())
	a.rootCmd.AddCommand(a.newDependenciesCommand())
	a.rootCmd.AddCommand(a.newDependentsCommand())

	return a.rootCmd.Execute()
}

func (a *App) newGetCommand() *cobra.Command {
	flags := NewGetFlags(a.GlobalFlags)
	runner := NewGetRunner(flags, a.IO)
	command := &cobra.Command{
		Use:   "get",
		Short: "Run terraform get",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run(cmd.Context()) },
	}
	return command
}

func (a *App) newDependenciesCommand() *cobra.Command {
	flags := NewDependenciesFlags(a.GlobalFlags)
	runner := NewDependenciesRunner(flags, terraform.NewDependencyStore(), a.IO)
	command := &cobra.Command{
		Use:   "dependencies",
		Short: "List module dependencies",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run() },
	}
	command.PersistentFlags().StringSliceVarP(&flags.StateDirs, "state-dirs", "s", []string{}, "state dirs")
	return command
}

func (a *App) newDependentsCommand() *cobra.Command {
	flags := NewDependentsFlags(a.GlobalFlags)
	runner := NewDependentsRunner(flags, terraform.NewDependentStore(), a.IO)
	command := &cobra.Command{
		Use:   "dependents",
		Short: "List module dependents",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run() },
	}
	command.PersistentFlags().StringSliceVarP(&flags.ModuleDirs, "module", "m", []string{}, "File paths of the module sources")
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
	log.SetOutput(io.Discard)
	if a.GlobalFlags.debug {
		log.SetOutput(a.IO.OutWriter)
	}
	log.SetPrefix(fmt.Sprintf("[%s] ", a.Ldflags.Name))
	log.Printf("Start: args: %v, ldflags: %+v", args, a.Ldflags)
}

type GlobalFlags struct {
	base   string
	format string
	debug  bool
}

func (f *GlobalFlags) GetBaseDir() *dir.BaseDir {
	return dir.NewBaseDir(f.base)
}

func (f *GlobalFlags) Format() string {
	return f.format
}

func (f *GlobalFlags) isDebug() bool {
	return f.debug
}

func (f *GlobalFlags) GoString() string {
	return fmt.Sprintf("{base: '%s', format: %s, debug: %t}", f.base, f.format, f.debug)
}
