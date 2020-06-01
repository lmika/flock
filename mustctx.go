package fabric

type MustContext struct {
	subContext *SubContext
}

type errMustContextFail struct {
	cause error
}

// RunEcho will run the command while printing stdout and stderr.
func (this *MustContext) RunEcho(cmd string, args ...string) {
	err := this.subContext.RunEcho(cmd, args...)
	if err != nil {
		panic(errMustContextFail{err})
	}
}

func (this *MustContext) Sudo() *MustContext {
	return &MustContext{
		subContext: &SubContext{
			driver:     this.subContext.driver,
			cmdBuilder: sudoCommandBuilder{this.subContext.cmdBuilder},
		},
	}
}
