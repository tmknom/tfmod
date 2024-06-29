package main

import (
	"log"
	"os"

	"github.com/tmknom/tfmod"
)

// Specify explicitly in ldflags
// For full details, see Makefile and .goreleaser.yml
var (
	name    = ""
	version = ""
	commit  = ""
	date    = ""
)

func main() {
	app := createApp()
	if err := app.Run(os.Args[1:]); err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func createApp() *tfmod.App {
	io := &tfmod.IO{
		InReader:  os.Stdin,
		OutWriter: os.Stdout,
		ErrWriter: os.Stderr,
	}

	ldflags := &tfmod.Ldflags{
		Name:    name,
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	return tfmod.NewApp(io, ldflags)
}
