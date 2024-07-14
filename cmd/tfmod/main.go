package main

import (
	"context"
	"log"
	"os"

	cli "github.com/tmknom/tfmod"
)

// Specify explicitly in ldflags
// For full details, see Makefile and .goreleaser.yml
var (
	name    = ""
	version = ""
	commit  = ""
	date    = ""
	url     = ""
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	io := &cli.IO{
		InReader:  os.Stdin,
		OutWriter: os.Stdout,
		ErrWriter: os.Stderr,
	}

	ldflags := &cli.Ldflags{
		Name:    name,
		Version: version,
		Commit:  commit,
		Date:    date,
		Url:     url,
	}

	ctx := context.Background()

	return cli.NewApp(io, ldflags).Run(ctx, os.Args[1:])
}
