package main

import (
	"context"
	"os/exec"
	"time"
)

// command will for most cases be used to mock external commands during testing.
var command = exec.CommandContext

type timeoutStep struct {
	step
	timeout time.Duration
}

func newTimeoutStep(name, cmd, message, projPath string, args []string, timeout time.Duration) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, cmd, message, projPath, args)
	s.timeout = timeout

	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}

	return s
}

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := command(ctx, s.cmd, s.args...)
	cmd.Dir = s.projPath

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed: timed out",
				cause: context.DeadlineExceeded,
			}
		}

		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}