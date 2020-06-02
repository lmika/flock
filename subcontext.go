package fabric

import (
	"context"
	"io"
	"io/ioutil"
)

type SubContext struct {
	driver     driver
	cmdBuilder commandBuilder
	fileDriver fileDriver
}

// RunEcho runs the commands with no input and all output going to stdout
func (sc *SubContext) RunEcho(cmd string, args ...string) error {
	builtCommand, err := sc.cmdBuilder.build(cmd, args)
	if err != nil {
		return err
	}

	return sc.driver.RunEcho(context.Background(), builtCommand)
}

func (this *SubContext) Sudo() *SubContext {
	newCmdBuilder := sudoCommandBuilder{this.cmdBuilder}

	return &SubContext{
		driver:     this.driver,
		cmdBuilder: newCmdBuilder,
		fileDriver: catFileDriver{this.driver, newCmdBuilder},
	}
}

// Run the next command
func (sc *SubContext) Must() *MustContext {
	return &MustContext{subContext: sc}
}

// MustDo will run the function in a "must" context (TODO: rename).  The commands in the must
// context must pass or the function will panic.  Panicing will return an error.
func (sc *SubContext) MustDo(fn func(*MustContext)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if failure, isFailure := e.(errMustContextFail); isFailure {
				err = failure.cause
			} else {
				panic(e)
			}
		}
	}()

	fn(&MustContext{subContext: sc})
	return err
}

// Open opens a remote file for reading.
func (sc *SubContext) Open(file string) (io.ReadCloser, error) {
	return sc.fileDriver.open(file)
}

// Create opens a remote file for writing.
func (sc *SubContext) Create(file string) (io.WriteCloser, error) {
	return sc.fileDriver.create(file)
}

// ReadFile reads the contents of a remote file.
func (sc *SubContext) ReadFile(file string) ([]byte, error) {
	r, err := sc.fileDriver.open(file)
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

// WriteFile writes the contents of a remote file.
func (sc *SubContext) WriteFile(file string, data []byte) error {
	w, err := sc.fileDriver.create(file)
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

// InDir returns a new subcontext for the current directory
func (sc *SubContext) InDir(newDir string) *SubContext {
	panic("TODO")
}
