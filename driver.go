package fabric

import (
	"context"
)

type driver interface {
	Close() error

	// Run a dokku command expecting no output
	Run(ctx context.Context, cmd string, args ...string) ([]byte, error)

	// Start starts a dokku command which will stream the output.  Onc
	RunEcho(ctx context.Context, cmd string, args ...string) error
}

