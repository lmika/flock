package flock

import (
	"github.com/melbahja/goph"
)

// SSHAuth is the means of authenticating a session.
type SSHAuth interface {
	toGophAuth() (goph.Auth, error)
}

type sshAuthFn func() (goph.Auth, error)

func (fn sshAuthFn) toGophAuth() (goph.Auth, error) {
	return fn()
}

// KeyPairAuth authenticates the session using a public-private key-pair, with the private key located at path with
// passphrase.
func KeyPairAuth(path string, passphrase string) SSHAuth {
	return sshAuthFn(func() (goph.Auth, error) {
		return goph.Key(path, passphrase), nil
	})
}

// PasswordAuth authenticates the session using a simple password.
func PasswordAuth(password string) SSHAuth {
	return sshAuthFn(func() (goph.Auth, error) {
		return goph.Password(password), nil
	})
}