package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
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

	pipeline := make([]step, 1)

	// CI-Step1: check if the app can build successfully.
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: Successful",
		projectPath,
		[]string{"build", ".", "errors"},
	)

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
