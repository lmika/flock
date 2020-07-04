package flock

import (
	"github.com/melbahja/goph"
)

type SSHAuth interface {
	toGophAuth() (goph.Auth, error)
}

type sshAuthFn func() (goph.Auth, error)

func (fn sshAuthFn) toGophAuth() (goph.Auth, error) {
	return fn()
}

func SSHKey(path string, passphrase string) SSHAuth {
	return sshAuthFn(func() (goph.Auth, error) {
		return goph.Key(path, passphrase), nil
	})
}
