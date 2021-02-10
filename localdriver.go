package flock

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
)

type localDriver struct {
	tracer	driverTracer
}

func (ld localDriver) Close() error {
	return nil
}

func (ld localDriver) Start(ctx context.Context, cmd *command, opts startOpts) (*runningCommand, error) {
	var err error
	osCmd := exec.CommandContext(ctx, cmd.name, cmd.args...)

	newRunningCommand := &runningCommand{
		ctx: ctx,
		command: cmd,
		tracer: ld.tracer,
	}

	if opts.pipeStdin {
		if newRunningCommand.stdin, err = osCmd.StdinPipe(); err != nil {
			return nil, err
		}
	}
	if opts.pipeStdout {
		if newRunningCommand.stdout, err = osCmd.StdoutPipe(); err != nil {
			return nil, err
		}
	}
	if opts.pipeStderr {
		if newRunningCommand.stderr, err = osCmd.StderrPipe(); err != nil {
			return nil, err
		}
	}

	if err = osCmd.Start(); err != nil {
		return nil, err
	}
	return newRunningCommand, err
}

func (ld localDriver) RunEcho(ctx context.Context, cmd *command) error {
	runningCmd, err := ld.Start(ctx, cmd, startOpts{
		pipeStdout: true,
		pipeStderr: true,
	})
	if err != nil {
		return err
	}

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go ld.echoPumpConsumer(runningCmd.stdout, waitGroup, cmd, false)
	go ld.echoPumpConsumer(runningCmd.stderr, waitGroup, cmd, true)
	waitGroup.Wait()

	return runningCmd.Wait()
}

func (ld localDriver) echoPumpConsumer(r io.Reader, wg *sync.WaitGroup, cmd *command, isStdErr bool) {
	defer wg.Done()

	scan := bufio.NewScanner(bufio.NewReader(r))
	for scan.Scan() {
		ld.tracer.echoOut(cmd, scan.Text(), isStdErr)
	}
}


type localFileDriver struct {}

// open opens a file for reading
func (localFileDriver) open(file string) (io.ReadCloser, error) {
	return os.Open(file)
}

// create opens a file for writing
func (localFileDriver) create(file string) (io.WriteCloser, error) {
	return os.Create(file)
}

// openAppend opens a file for appending
func (localFileDriver) openAppend(file string) (io.WriteCloser, error) {
	return os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}