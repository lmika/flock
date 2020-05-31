package fabric

import (
	"bufio"
	"context"
	"fmt"
	"github.com/melbahja/goph"
	"github.com/pkg/errors"
	"io"
	"log"
	"strings"
	"sync"
)

type sshDriver struct {
	client *goph.Client
}

func (this *sshDriver) Close() error {
	return this.client.Close()
}

func (this *sshDriver) Run(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	// TODO: escape args
	fullCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))

	output, err := this.client.Run(fullCmd)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot run "+cmd)
	}

	return output, nil
}


func (this *sshDriver) RunEcho(ctx context.Context, cmd string, args ...string) error {
	fullCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))

	session, err := this.client.NewSession()
	if err != nil {
		return errors.Wrapf(err, "cannot start session")
	}

	stdoutReader, stdoutWriter := io.Pipe()
	session.Stdout = stdoutWriter

	stderrReader, stderrWriter := io.Pipe()
	session.Stderr = stderrWriter

	if err := session.Start(fullCmd); err != nil {

		return errors.Wrapf(err, "cannot start command: %v", fullCmd)
	}

	outChan, errChan, doneChan := make(chan string), make(chan string), make(chan error)

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)
	go pumpToStreamEvent(outChan, stdoutReader, waitGroup)
	go pumpToStreamEvent(errChan, stderrReader, waitGroup)
	go func() {doneChan <- session.Wait()}()

	for {
		select {
		case err := <-doneChan:
			session.Close()
			stdoutWriter.Close()
			stderrWriter.Close()
			waitGroup.Wait()
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func pumpToStreamEvent(dest chan string, r io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()

	scan := bufio.NewScanner(bufio.NewReader(r))
	for scan.Scan() {
		// TODO: instead of logging, allow another component to receive echo commands
		log.Println("[remote] ", scan.Text())
	}
}



