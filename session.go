package fabric

import (
	"github.com/melbahja/goph"
	"log"
	"os"
)

type Session struct {
	*SubContext
}

func NewSSH(user, addr string, auth SSHAuth) (*Session, error) {
	gophAuth, err := auth.toGophAuth()
	if err != nil {
		return nil, err
	}

	//sshClient, err := goph.New(user, addr, gophAuth)

	// TODO: make this configurable
	sshClient, err := goph.NewUnknown(user, addr, gophAuth)
	if err != nil {
		return nil, err
	}

	// TODO: make THIS configurable
	logTag := user + "@" + addr

	driver := &sshDriver{
		client: sshClient,
		// TODO: make configurable
		tracer: &noisyLogDriverTracer{
			logger: log.New(os.Stderr, "", log.Ldate | log.Ltime),
			prefix: logTag,
		},
	}
	return &Session{
		&SubContext{
			driver:     driver,
			cmdBuilder: plainCommandBuilder{},
		},
	}, nil
}

func (this *Session) Close() error {
	return this.driver.Close()
}

func (this *Session) Sudo() *SubContext {
	return &SubContext{
		driver:     this.driver,
		cmdBuilder: sudoCommandBuilder{this.cmdBuilder},
	}
}
