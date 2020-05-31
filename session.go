package fabric

import (
	"github.com/melbahja/goph"
	"os"
)

type Session struct {
	driver	driver
}

func NewSSH(user, addr string, auth goph.Auth) (*Session, error) {
	sshClient, err := goph.New("root", "lmika.app", goph.Key(os.ExpandEnv("${HOME}/.ssh/id_rsa"), ""))
	if err != nil {
		return nil, err
	}
	return &Session{&sshDriver{sshClient}}, nil
}

func (this *Session) Close() error {
	return this.driver.Close()
}

// MustDo will run the function in a "must" context (TODO: rename).  The commands in the must
// context must pass or the function will panic.  Panicing will return an error.
func (this *Session) MustDo(fn func(*MustContext)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if failure, isFailure := e.(errMustContextFail); isFailure {
				err = failure.cause
			} else {
				panic(e)
			}
		}
	}()

	fn(&MustContext{this.driver})
	return err
}
