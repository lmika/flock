package fabric

import (
	"io"
)

type MustContext struct {
	subContext *SubContext
}

type errMustContextFail struct {
	cause error
}

// Panic will cause the current MustContext to panic and fail.  When running in MustDo() closure,
// the passed in error will be returned.
func (this *MustContext) Panic(err error) {
	panic(errMustContextFail{err})
}

// RunEcho will run the command while printing stdout and stderr.
func (this *MustContext) RunEcho(cmd string, args ...string) {
	err := this.subContext.RunEcho(cmd, args...)
	if err != nil {
		this.Panic(err)
	}
}

// Sudo returns a new sudo MustContext.
func (this *MustContext) Sudo() *MustContext {
	return &MustContext{
		subContext: this.subContext.Sudo(),
	}
}

// Open opens a remote file for reading.  This will panic if there is a problem opening the file.  Any errors
// reading or closing the file will not result in a panic.
func (this *MustContext) Open(file string) io.ReadCloser {
	rc, err := this.subContext.Open(file)
	if err != nil {
		this.Panic(err)
	}
	return rc
}

// Create opens a remote file for writing.  This will panic if there is a problem creating the remote file.
// An errors writing or closing to the file will not result in a panic.
func (this *MustContext) Create(file string) io.WriteCloser {
	wc, err := this.subContext.Create(file)
	if err != nil {
		this.Panic(err)
	}
	return wc
}

// ReadFile reads the contents of a remote file.  This will panic if there was a problem opening, reading or
// closing the file.
func (this *MustContext) ReadFile(file string) []byte {
	bts, err := this.subContext.ReadFile(file)
	if err != nil {
		this.Panic(err)
	}
	return bts
}

// WriteFile writes the contents of a remote file.  This will panic if there was a problem opening, writing to,
// or closing the file.
func (this *MustContext) WriteFile(file string, data []byte) {
	err := this.subContext.WriteFile(file, data)
	if err != nil {
		this.Panic(err)
	}
}