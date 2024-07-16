package tfmod

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

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

func (i *IO) Read() []string {
	lines := make([]string, 0, 64)
	scanner := bufio.NewScanner(i.InReader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (i *IO) IsPipe() bool {
	stat, err := i.InReader.(fs.File).Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

type Ldflags struct {
	Name    string
	Version string
	Commit  string
	Date    string
	Url     string
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

func (a *App) Run(ctx context.Context, args []string) error {
	a.prepareCommand(ctx, args)

	a.rootCmd.PersistentFlags().StringVarP(&a.GlobalFlags.base, "base", "b", ".", "The base directory that contains tf files")
	a.rootCmd.PersistentFlags().StringVarP(&a.GlobalFlags.format, "format", "f", format.TextFormat, fmt.Sprintf("Format output by: {%s}", format.SupportType()))
	a.rootCmd.PersistentFlags().BoolVar(&a.GlobalFlags.debug, "debug", false, "Show debugging output")

	a.rootCmd.AddCommand(a.newDownloadCommand())
	a.rootCmd.AddCommand(a.newDependencyCommand())
	a.rootCmd.AddCommand(a.newDependentCommand())

	return a.rootCmd.Execute()
}

func (a *App) newDownloadCommand() *cobra.Command {
	flags := NewDownloadFlags(a.GlobalFlags)
	runner := NewDownloadRunner(flags, a.IO)
	command := &cobra.Command{
		Use:   "download",
		Short: "Download all modules at once under the base directory",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run(cmd.Context()) },
	}
	command.PersistentFlags().IntVar(&flags.MaxConcurrent, "max-concurrent", defaultMaxConcurrent, "Maximum number of concurrency for terraform get")
	return command
}

func (a *App) newDependencyCommand() *cobra.Command {
	flags := NewDependencyFlags(a.GlobalFlags)
	runner := NewDependencyRunner(flags, terraform.NewDependencyStore(), a.IO)
	command := &cobra.Command{
		Use:   "dependency",
		Short: "Explore the module directory it depends on",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run() },
	}
	command.PersistentFlags().StringSliceVarP(&flags.StatePaths, "state", "s", []string{}, "Directory paths for the state to managed configuration")
	return command
}

func (a *App) newDependentCommand() *cobra.Command {
	flags := NewDependentFlags(a.GlobalFlags)
	runner := NewDependentRunner(flags, terraform.NewDependentStore(), a.IO)
	command := &cobra.Command{
		Use:   "dependent",
		Short: "Explore how the state directory is used by specified modules",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.Run() },
	}
	command.PersistentFlags().StringSliceVarP(&flags.ModulePaths, "module", "m", []string{}, "File paths of the module sources")
	return command
}

func (a *App) prepareCommand(ctx context.Context, args []string) {
	a.rootCmd.SetContext(ctx)

	// setup help message
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
	version := fmt.Sprintf("%s version %s (%s)\n%s\n", a.Ldflags.Name, a.Ldflags.Version, a.Ldflags.Date, a.Ldflags.Url)
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

func (f *GlobalFlags) BaseDir() *dir.BaseDir {
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
