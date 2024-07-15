package terraform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/errlib"
)

type Command struct {
	maxConcurrent int
}

func NewCommand(maxConcurrent int) *Command {
	return &Command{
		maxConcurrent: maxConcurrent,
	}
}

type token struct{}

func (c *Command) GetAll(ctx context.Context, workDirs []*dir.Dir) error {
	var wg sync.WaitGroup
	wg.Add(len(workDirs))

	sem := make(chan token, c.maxConcurrent)
	resultCh := make(chan string, len(workDirs))
	errCh := make(chan error, len(workDirs))
	for _, arg := range workDirs {
		sem <- token{}
		go func(arg string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			info, err := c.executeGet(ctx, arg)
			if err != nil {
				errCh <- err
			} else {
				resultCh <- info
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

func (c *Command) executeGet(ctx context.Context, workDir string) (string, error) {
	cmd := exec.CommandContext(ctx, "terraform", "get")
	cmd.Dir = workDir
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	err := cmd.Run()
	if err != nil {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%s\n", cmd.String()))
		b.WriteString(fmt.Sprintf("Stderr\n%s\n", cmd.Stderr.(*bytes.Buffer).String()))
		b.WriteString(fmt.Sprintf("Stdout\n%s\n", cmd.Stdout.(*bytes.Buffer).String()))
		b.WriteString(fmt.Sprintf("Workdir: %v\n", cmd.Dir))
		return "", errlib.Wrapf(err, "%s", b.String())
	}
	return fmt.Sprintf("execute: %s (at %s)", cmd.String(), cmd.Dir), nil
}
