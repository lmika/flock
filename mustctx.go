package fabric

import "context"

type MustContext struct {
	driver		driver
}

type errMustContextFail struct {
	cause error
}

// RunEcho will run the command while printing stdout and stderr.
func (this *MustContext) RunEcho(cmd string, args ...string) {
	err := this.driver.RunEcho(context.Background(), cmd, args...)
	if err != nil {
		panic(errMustContextFail{err})
	}
}
