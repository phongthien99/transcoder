package model

import (
	"fmt"
	"io"
	"os/exec"
)

type Transcoder struct {
	StdErrPipe   io.ReadCloser
	StdStdinPipe io.WriteCloser
	Cmd          *exec.Cmd
	InputInfo    Metadata
	Progress     Progress
}

func (t *Transcoder) GetInputInfo(input string) (*Metadata, error) {
	return nil, nil
}

func (t *Transcoder) GetProcess() (*Progress, error) {
	return nil, nil
}

func (t *Transcoder) Run(progress bool) <-chan error {

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		stderr, err := t.Cmd.StderrPipe()
		if err != nil {
			errChan <- err
			return
		}
		t.StdErrPipe = stderr
		stdin, err := t.Cmd.StdinPipe()
		if err != nil {
			errChan <- err
			return
		}
		t.StdStdinPipe = stdin
		if err := t.Cmd.Start(); err != nil {
			errChan <- err
			return
		}

		if progress {
		}

		if err := t.Cmd.Wait(); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()
	return errChan
}
func (t *Transcoder) Stop() error {
	if t.Cmd != nil && t.Cmd.Process != nil {
		return t.Cmd.Process.Kill()
	}
	return fmt.Errorf("no running process to stop")
}
