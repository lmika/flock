package fabric

import (
	"fmt"
	"strings"
)

type ttyMode int

const (
	ttyDontCare ttyMode = iota
	ttyNever
	ttyRecommended
	ttyRequired
)

type command struct {
	name string
	args []string
	ttyMode ttyMode
}

func (cmd *command) fullCommand() string {
	return fmt.Sprintf("%s %s", cmd.name, strings.Join(cmd.args, " "))
}


type commandBuilder interface {
	build(name string, args []string) (*command, error)
}


type plainCommandBuilder struct{}

func (plainCommandBuilder) build(name string, args []string) (*command, error) {
	return &command{
		name: name,
		args: args,
	}, nil
}


type sudoCommandBuilder struct {
	delegate commandBuilder
}

func (b sudoCommandBuilder) build(name string, args []string) (*command, error) {
	delegatedCmd, err := b.delegate.build(name, args)
	if err != nil {
		return nil, err
	}

	return &command{
		name: "sudo",
		args: append([]string{delegatedCmd.name}, delegatedCmd.args...),
		ttyMode: ttyRecommended,
	}, nil
}