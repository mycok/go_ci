package main

import (
	"os/exec"
)

type step struct {
	name     string
	cmd      string // name of the executable command
	args     []string
	message  string // output message on successful execution
	projPath string // path / directory on which to execute the step
}

func newStep(name, cmd, message, projPath string, args []string) step {
	return step{
		name:     name,
		cmd:      cmd,
		args:     args,
		message:  message,
		projPath: projPath,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.cmd, s.args...)
	cmd.Dir = s.projPath

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "go build failed",
			cause: err,
		}
	}

	return s.message, nil
}
