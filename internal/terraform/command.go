package terraform

import (
	"bytes"
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

func (c *Command) GetAll(workDirs []*dir.Dir) error {
	jobs := make(chan string, len(workDirs))
	var wg sync.WaitGroup

	for i := 0; i < maxConcurrentJobs; i++ {
		wg.Add(1)
		go c.worker(jobs, &wg)
	}

	for _, workDir := range workDirs {
		jobs <- workDir.Abs()
	}
	close(jobs)

	wg.Wait()

	return nil
}

func (c *Command) worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for workDir := range jobs {
		err := c.executeGet(workDir)
		if err != nil {
			log.Printf("Error terraform get in %s: %v\n", workDir, err)
		}
	}
}

func (c *Command) executeGet(workDir string) error {
	cmd := exec.Command("terraform", "get")
	cmd.Dir = workDir
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	info := fmt.Sprintf("%s (at %s)\n", cmd.String(), cmd.Dir)
	log.Printf(fmt.Sprintf("execute: %s", info))

	err := cmd.Run()
	if err != nil {
		return errlib.Wrapf(err, "%s\n stdout: %v\n stderr: %v\n", info, cmd.Stdout, cmd.Stderr)
	}
	return nil
}
