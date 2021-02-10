package flock

import (
	"io"
)

// MustContext represents the execution environment, such as the current directory, user or machine.
// It's a lot like context except that any failures will cause the MustContext to panic instead of
// return an error.
type MustContext struct {
	subContext *Context
}

type errMustContextFail struct {
	cause error
}

// Panic will cause the current MustContext to panic and fail.  When running in Context.MustDo(),
// the passed in error will be returned.
func (mustCtx *MustContext) Panic(err error) {
	panic(errMustContextFail{err})
}

// RunEcho runs the commands with no input and all output echoed to stdout.  The context will panic
// if the command fails to execute successfully.
func (mustCtx *MustContext) RunEcho(cmd string, args ...string) {
	err := mustCtx.subContext.RunEcho(cmd, args...)
	if err != nil {
		mustCtx.Panic(err)
	}
}

// Sudo returns a new context which will run commands and transfer files as the sudo user.
func (mustCtx *MustContext) Sudo() *MustContext {
	return &MustContext{
		subContext: mustCtx.subContext.Sudo(),
	}
}

// Open opens a remote file for reading.  The context will panic if there is a problem opening the file.  Any errors
// reading or closing the file will not result in a panic.  The caller must close the file once done.
func (mustCtx *MustContext) Open(file string) io.ReadCloser {
	rc, err := mustCtx.subContext.Open(file)
	if err != nil {
		mustCtx.Panic(err)
	}
	return rc
}

// Create opens a remote file for writing.  The context will panic if there is a problem creating the remote file.
// An errors writing or closing to the file will not result in a panic.  The caller must close the file once done.
func (mustCtx *MustContext) Create(file string) io.WriteCloser {
	wc, err := mustCtx.subContext.Create(file)
	if err != nil {
		mustCtx.Panic(err)
	}
	return wc
}

// ReadFile reads the contents of a remote file.  The context will panic if there was a problem opening, reading or
// closing the file.
func (mustCtx *MustContext) ReadFile(file string) []byte {
	bts, err := mustCtx.subContext.ReadFile(file)
	if err != nil {
		mustCtx.Panic(err)
	}
	return bts
}

// WriteFile writes the contents of a remote file.  The context will panic if there was a problem opening, writing to,
// or closing the file.
func (mustCtx *MustContext) WriteFile(file string, data []byte) {
	err := mustCtx.subContext.WriteFile(file, data)
	if err != nil {
		mustCtx.Panic(err)
	}
}

// AppendToFile adds the contents to the end of a remote file.  The context will panic if there was a problem opening, writing to,
// or closing the file.
func (mustCtx *MustContext) AppendToFile(file string, data []byte) {
	err := mustCtx.subContext.AppendToFile(file, data)
	if err != nil {
		mustCtx.Panic(err)
	}
}

// Upload copies the contents of a local file to the remote machine.  The context will panic if there was a problem
// within the process.
func (mustCtx *MustContext) Upload(remoteFile, localFile string) {
	err := mustCtx.subContext.Upload(remoteFile, localFile)
	if err != nil {
		mustCtx.Panic(err)
	}
}

// Download copies the contents of a remote file to a local file.  The context will panic if there was a problem
// within the process.
func (mustCtx *MustContext) Download(localFile, remoteFile string) {
	err := mustCtx.subContext.Download(localFile, remoteFile)
	if err != nil {
		mustCtx.Panic(err)
	}
}