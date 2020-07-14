package fabric

import (
	"context"
	"io"
)

type fileDriver interface {
	// open opens a file for reading
	open(file string) (io.ReadCloser, error)

	// create opens a file for writing
	create(file string) (io.WriteCloser, error)

	// openAppend opens a file for appending
	openAppend(file string) (io.WriteCloser, error)
}

// fileDriver that uses 'cat' to read/write files
type catFileDriver struct {
	driver     driver
	cmdBuilder commandBuilder
}

func (cfd catFileDriver) open(file string) (io.ReadCloser, error) {
	cmd, err := cfd.cmdBuilder.build("cat", []string{file})
	if err != nil {
		return nil, err
	}

	runningCmd, err := cfd.driver.Start(context.Background(), cmd, startOpts{pipeStdout: true, ttyMode: ttyNever})
	if err != nil {
		return nil, err
	}

	return runningCmdReaderCloser{runningCmd, runningCmd.stdout}, nil
}

func (cfd catFileDriver) create(file string) (io.WriteCloser, error) {
	// TODO: properly escape this command!!!
	cmd, err := cfd.cmdBuilder.build("sh", []string{"-c", "'cat > " + file + "'"})
	if err != nil {
		return nil, err
	}

	// TODO: capture output in stderr
	runningCmd, err := cfd.driver.Start(context.Background(), cmd, startOpts{pipeStdin: true, ttyMode: ttyNever})
	if err != nil {
		return nil, err
	}

	return runningCmdWriterCloser{runningCmd, runningCmd.stdin}, nil
}

func (cfd catFileDriver) openAppend(file string) (io.WriteCloser, error) {
	// TODO: properly escape this command!!!
	cmd, err := cfd.cmdBuilder.build("sh", []string{"-c", "'cat >> " + file + "'"})
	if err != nil {
		return nil, err
	}

	// TODO: capture output in stderr
	runningCmd, err := cfd.driver.Start(context.Background(), cmd, startOpts{pipeStdin: true, ttyMode: ttyNever})
	if err != nil {
		return nil, err
	}

	return runningCmdWriterCloser{runningCmd, runningCmd.stdin}, nil
}

type runningCmdReaderCloser struct {
	cmd *runningCommand
	r   io.Reader
}

func (r runningCmdReaderCloser) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r runningCmdReaderCloser) Close() error {
	return r.cmd.Wait()
}

type runningCmdWriterCloser struct {
	cmd *runningCommand
	w   io.WriteCloser
}

func (w runningCmdWriterCloser) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w runningCmdWriterCloser) Close() error {
	if err := w.w.Close(); err != nil {
		return err
	}
	return w.cmd.Wait()
}
