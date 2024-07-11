package terraform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/errlib"
)

type Command struct{}

func NewCommand() *Command {
	return &Command{}
}

type token struct{}

func (c *Command) GetAll(ctx context.Context, workDirs []*dir.Dir) error {
	var wg sync.WaitGroup
	wg.Add(len(workDirs))

	sem := make(chan token, maxConcurrentJobs)
	resultCh := make(chan string, len(workDirs))
	errCh := make(chan error, len(workDirs))
	for _, arg := range workDirs {
		sem <- token{}
		go func(arg string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			err := c.executeGet(ctx, arg)
			if err != nil {
				errCh <- err
			}
		}(arg.Abs())
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(resultCh)
		close(sem)
	}()

	for result := range resultCh {
		log.Println(result)
	}

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (c *Command) executeGet(ctx context.Context, workDir string) error {
	cmd := exec.CommandContext(ctx, "terraform", "get")
	cmd.Dir = workDir
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	info := fmt.Sprintf("%s (at %s)\n", cmd.String(), cmd.Dir)
	log.Printf(fmt.Sprintf("execute: %s", info))

	err := cmd.Run()
	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		return errlib.Wrapf(err, "%s\nStderr:\n%v\nStdout:\n%v\nWorkdir: %v\nExitcode: %d\n", cmd.String(), cmd.Stderr.(*bytes.Buffer).String(), cmd.Stdout.(*bytes.Buffer).String(), cmd.Dir, exitCode)
	}
	return nil
}
