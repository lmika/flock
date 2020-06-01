package fabric

import (
	"bufio"
	"context"
	"github.com/melbahja/goph"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
)

type sshDriver struct {
	client *goph.Client
	tracer	driverTracer
}

func (this *sshDriver) Close() error {
	return this.client.Close()
}

func (this *sshDriver) Start(ctx context.Context, cmd *command, opts startOpts) (*runningCommand, error) {
	fullCmd := cmd.fullCommand()

	session, err := this.client.NewSession()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot start session")
	}

	if err := this.requestTtyIfRequired(session, cmd, opts); err != nil {
		return nil, errors.Wrapf(err, "cannot start session")
	}

	newRunningCommand := &runningCommand{
		ctx: ctx,
		command: cmd,
		tracer: this.tracer,
		session: session,
	}

	if opts.pipeStdin {
		if newRunningCommand.stdin, err = session.StdinPipe(); err != nil {
			return nil, err
		}
	}
	if opts.pipeStdout {
		if newRunningCommand.stdout, err = session.StdoutPipe(); err != nil {
			return nil, err
		}
	}
	if opts.pipeStderr {
		if newRunningCommand.stderr, err = session.StderrPipe(); err != nil {
			return nil, err
		}
	}

	this.tracer.startCmd(cmd)
	if err := session.Start(fullCmd); err != nil {
		return nil, errors.Wrapf(err, "cannot start command: %v", fullCmd)
	}

	return newRunningCommand, nil
}


func (this *sshDriver) Run(ctx context.Context, cmd *command) ([]byte, error) {
	// TODO: escape args
	fullCmd := cmd.fullCommand()

	this.tracer.startCmd(cmd)
	output, err := this.client.Run(fullCmd)
	this.tracer.endCommand(cmd, err, false)

	if err != nil {
		return nil, errors.Wrapf(err, "cannot run "+cmd.name)
	}

	return output, nil
}

func (this *sshDriver) RunEcho(ctx context.Context, cmd *command) error {
	fullCmd := cmd.fullCommand()

	session, err := this.client.NewSession()
	if err != nil {
		return errors.Wrapf(err, "cannot start session")
	}

	if err := this.requestTtyIfRequired(session, cmd, startOpts{}); err != nil {
		return errors.Wrapf(err, "cannot start session")
	}

	stdoutReader, stdoutWriter := io.Pipe()
	session.Stdout = stdoutWriter

	stderrReader, stderrWriter := io.Pipe()
	session.Stderr = stderrWriter

	this.tracer.startCmd(cmd)
	if err := session.Start(fullCmd); err != nil {
		return errors.Wrapf(err, "cannot start command: %v", fullCmd)
	}

	doneChan := make(chan error)

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)
	go this.echoPumpConsumer(stdoutReader, waitGroup, cmd, false)
	go this.echoPumpConsumer(stderrReader, waitGroup, cmd, true)
	go func() {doneChan <- session.Wait()}()

	for {
		select {
		case err := <-doneChan:
			session.Close()
			stdoutWriter.Close()
			stderrWriter.Close()
			waitGroup.Wait()

			this.tracer.endCommand(cmd, err, false)
			return err
		case <-ctx.Done():
			this.tracer.endCommand(cmd, ctx.Err(), true)
			return ctx.Err()
		}
	}
}

func (this *sshDriver) echoPumpConsumer(r io.Reader, wg *sync.WaitGroup, cmd *command, isStdErr bool) {
	defer wg.Done()

	scan := bufio.NewScanner(bufio.NewReader(r))
	for scan.Scan() {
		this.tracer.echoOut(cmd, scan.Text(), isStdErr)
	}
}

func (this *sshDriver) requestTtyIfRequired(session *ssh.Session, cmd *command, opts startOpts) error {
	// TODO: make this configurable
	effectiveTtyMode := ttyNever

	if cmd.ttyMode != ttyDontCare {
		effectiveTtyMode = cmd.ttyMode
	}

	if opts.ttyMode != ttyDontCare {
		effectiveTtyMode = opts.ttyMode
	}

	if effectiveTtyMode == ttyNever {
		return nil
	}

	if err := session.RequestPty("xterm", 40, 80, ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}); err != nil {
		if effectiveTtyMode == ttyRequired {
			return errors.Wrapf(err, "cannot start tty")
		}
		this.tracer.warnf("unable to obtain recommended tty: %v", err)
	}

	return nil
}