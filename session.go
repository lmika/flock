package flock

import (
	"github.com/melbahja/goph"
	"log"
	"os"
)

// Session represents an SSH session to a particular machine.
// A session must be closed once it is no longer required.
type Session struct {
	*Context
}

// LocalSession returns a session that is configured to run on the local machine.  The commands and file transfer
// operations are run directly on the machine itself: they do not require a running SSH service.
func LocalSession() *Session {
	return &Session{
		&Context{
			driver:     localDriver{
				tracer: &noisyLogDriverTracer{
					logger: log.New(os.Stderr, "", log.Ldate | log.Ltime),
					prefix: "local",
				},
			},
			cmdBuilder: plainCommandBuilder{},

			// TODO: the default session should use the SSH connection sessions here to save having to use 'cat'
			fileDriver: localFileDriver{},
		},
	}
}

// NewSSH creates a new SSH session to the remote machine addr as user authenticated using auth.
// An error is returned if there is a problem establishing the SSH session.
func NewSSH(addr, user string, auth SSHAuth) (*Session, error) {
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
		&Context{
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

// Close closes the session.
func (session *Session) Close() error {
	return session.driver.Close()
}