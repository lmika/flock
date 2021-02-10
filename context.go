package flock

import (
	"context"
	"io"
	"io/ioutil"
	"os"
)

// Context represents the execution environment, such as the current directory, user, or machine.
// Contexts are immutable: any operation that will result in a change in the environment will return a new
// context.
type Context struct {
	driver     driver
	cmdBuilder commandBuilder
	fileDriver fileDriver
}

// RunEcho runs the commands with no input and all output echoed to stdout.  An error is returned
// if the command fails to execute successfully.
func (ctx *Context) RunEcho(cmd string, args ...string) error {
	builtCommand, err := ctx.cmdBuilder.build(cmd, args)
	if err != nil {
		return err
	}

	return ctx.driver.RunEcho(context.Background(), builtCommand)
}

// Sudo returns a new context which will run commands and transfer files as the sudo user.
func (ctx *Context) Sudo() *Context {
	newCmdBuilder := sudoCommandBuilder{ctx.cmdBuilder}

	return &Context{
		driver:     ctx.driver,
		cmdBuilder: newCmdBuilder,
		fileDriver: catFileDriver{ctx.driver, newCmdBuilder},
	}
}

// Must returns a must context backed by this context.
func (ctx *Context) Must() *MustContext {
	return &MustContext{subContext: ctx}
}

// MustDo will run the function with a must context that mirrors the current context.
// Any panics thrown by the MustContext as a result of a failed command will be recovered
// here and retured as an error.  Any other panics will be left to propogate up the stack.
func (ctx *Context) MustDo(fn func(*MustContext)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if failure, isFailure := e.(errMustContextFail); isFailure {
				err = failure.cause
			} else {
				panic(e)
			}
		}
	}()

	fn(&MustContext{subContext: ctx})
	return err
}

// Open opens a remote file for reading.  The caller must close the file once done.
func (ctx *Context) Open(file string) (io.ReadCloser, error) {
	return ctx.fileDriver.open(file)
}

// Create opens a remote file for writing.  The caller must close the file once done.
func (ctx *Context) Create(file string) (io.WriteCloser, error) {
	return ctx.fileDriver.create(file)
}

// ReadFile reads the contents of a remote file and returns it as a byte slice.
func (ctx *Context) ReadFile(file string) ([]byte, error) {
	r, err := ctx.fileDriver.open(file)
	if err != nil {
		return nil, err
	}

	bts, err := ioutil.ReadAll(r)
	if err != nil {
		r.Close()
		return nil, err
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	return bts, nil
}

// WriteFile writes the contents of a remote file.  If the file exists, it will be overwritten.
func (ctx *Context) WriteFile(file string, data []byte) error {
	w, err := ctx.fileDriver.create(file)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

// AppendToFile adds the contents to the end of a remote file.  If the file does not exist, it will be created.
func (ctx *Context) AppendToFile(file string, data []byte) error {
	w, err := ctx.fileDriver.openAppend(file)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

// Upload copies the contents of a local file to the remote machine.
func (ctx *Context) Upload(remoteFile, localFile string) error {
	f, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer f.Close()

	wc, err := ctx.Create(remoteFile)
	if err != nil {
		return err
	}

	if _, err = io.Copy(wc, f); err != nil {
		wc.Close()
		return err
	}

	return wc.Close()
}

// Download copies the contents of a remote file to a local file.
func (ctx *Context) Download(localFile, remoteFile string) error {
	rc, err := ctx.Open(remoteFile)
	if err != nil {
		return err
	}
	defer rc.Close()

	f, err := os.Create(localFile)
	if err != nil {
		return err
	}

	if _, err = io.Copy(f, rc); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}
