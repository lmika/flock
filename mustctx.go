package fabric

import "context"

type MustContext struct {
	driver		driver
	cmdBuilder	commandBuilder
}

type errMustContextFail struct {
	cause error
}

// RunEcho will run the command while printing stdout and stderr.
func (this *MustContext) RunEcho(cmd string, args ...string) {
	builtCmd, err := this.cmdBuilder.build(cmd, args)
	if err != nil {
		panic(errMustContextFail{err})
	}

	err = this.driver.RunEcho(context.Background(), builtCmd)
	if err != nil {
		panic(errMustContextFail{err})
	}
}
