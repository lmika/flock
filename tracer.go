package flock

import (
	"fmt"
	"log"
)

type driverTracer interface {
	warnf(msg string, args ...interface{})
	startCmd(cmd *command)
	endCommand(cmd *command, err error, cancelled bool)
	echoOut(cmd *command, line string, fromStderr bool)
}

type noisyLogDriverTracer struct {
	logger	*log.Logger
	prefix	string
}

func (lt *noisyLogDriverTracer) startCmd(cmd *command) {
	lt.logger.Printf("[%s] START %v", lt.prefix, cmd.fullCommand())
}

func (lt *noisyLogDriverTracer) warnf(msg string, args ...interface{}) {
	logLine := fmt.Sprintf(msg, args...)
	lt.logger.Printf("[%s] WARN  %v", lt.prefix, logLine)
}

func (lt *noisyLogDriverTracer) endCommand(cmd *command, err error, cancelled bool) {
	if err == nil {
		lt.logger.Printf("[%s] DONE  %v", lt.prefix, cmd.fullCommand())
	} else if cancelled {
		lt.logger.Printf("[%s] CANCL %v", lt.prefix, cmd.fullCommand())
	} else {
		lt.logger.Printf("[%s] ERR   %v", lt.prefix, cmd.fullCommand())
		lt.logger.Printf("[%s] ERR   .. %v", lt.prefix, err)
	}
}

func (lt *noisyLogDriverTracer) echoOut(cmd *command, line string, fromStderr bool) {
	if fromStderr {
		lt.logger.Printf("[%s] . ERR %v", lt.prefix, line)
	} else {
		lt.logger.Printf("[%s] . OUT %v", lt.prefix, line)
	}
}
