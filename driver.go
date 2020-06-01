package fabric

import (
	"context"
	"golang.org/x/crypto/ssh"
	"io"
)

type driver interface {
	Close() error

	Start(ctx context.Context, cmd *command, opts startOpts) (*runningCommand, error)

	// Run a dokku command expecting no output
	Run(ctx context.Context, cmd *command) ([]byte, error)

	// Start starts a dokku command which will stream the output.  Onc
	RunEcho(ctx context.Context, cmd *command) error
}

type startOpts struct {
	pipeStdout, pipeStdin, pipeStderr bool
}

type runningCommand struct {
	ctx     context.Context
	command *command
	tracer  driverTracer
	session *ssh.Session

	stdin  io.WriteCloser
	stdout io.Reader
	stderr io.Reader
}

func (rc *runningCommand) Wait() error {
	doneChan := make(chan error)
	go func() { doneChan <- rc.session.Wait() }()

	for {
		select {
		case err := <-doneChan:
			rc.session.Close()
			rc.tracer.endCommand(rc.command, err, false)
			return err
		case <-rc.ctx.Done():
			rc.tracer.endCommand(rc.command, rc.ctx.Err(), true)
			return rc.ctx.Err()
		}
	}
}
