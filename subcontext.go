package fabric

import "context"

type SubContext struct {
	driver     driver
	cmdBuilder commandBuilder
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
	return &SubContext{
		driver:     this.driver,
		cmdBuilder: sudoCommandBuilder{this.cmdBuilder},
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