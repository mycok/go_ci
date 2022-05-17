package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main()  {
	proj := flag.String("proj", "", "project path or directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(projectPath string, out io.Writer) error {
	if projectPath == "" {
		return fmt.Errorf("project path or directory is required: %w", ErrValidation)
	}

	// CI-Step1: check if the app can build successfully.
	cmdArgs := []string{"build", ".", "errors"}
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return &stepErr{step: "go build", msg: "go build failed", cause: err}
	}

	_, err := fmt.Fprintln(out, "Go Build: Successful")
	
	return err
}