package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// TODO: Add functionality to enable the execute method write all cmd out to
// the standard output. achieve this by adding a writer arg to the excute method.

type executionStep struct {
	step
}

func newExecutionStep(name, cmd, message, projPath string, args []string) executionStep {
	s := executionStep{}
	s.step = newStep(name, cmd, message, projPath, args)

	return s
}

// excute is an override of the original step.execute() method. this handles
// cases where the executed step command doesn't return expected errors on
// unsuccessful executions.
func (s executionStep) execute() (string, error) {
	var out bytes.Buffer

	cmd := exec.Command(s.cmd, s.args...)
	cmd.Stdout = &out
	cmd.Dir = s.projPath

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	if out.Len() > 0 {
		return "", &stepErr{
			step:  s.name,
			msg:   fmt.Sprintf("invalid format for file: %s", out.String()),
			cause: nil,
		}
	}

	return s.message, nil
}
