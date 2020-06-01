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
	cmdBuilder := plainCommandBuilder{}

	return &Session{
		&SubContext{
			driver:     driver,
			cmdBuilder: cmdBuilder,

			// TODO: the default session should use the SSH connection sessions here to save having to use 'cat'
			fileDriver: catFileDriver{
				driver: driver,
				cmdBuilder: cmdBuilder,
			},
		},
	}, nil
}

func (this *Session) Close() error {
	return this.driver.Close()
}

//func (this *Session) Tunnel() *Tunnel {
//
//}