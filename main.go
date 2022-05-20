package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	projPath := flag.String("proj", "", "project path or directory")
	flag.Parse()

	if err := run(*projPath, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(projPath string, out io.Writer) error {
	if projPath == "" {
		return fmt.Errorf("project path or directory is required: %w", ErrValidation)
	}

	pipeline := make([]executer, 4)

	// CI-Step1: check if the app can build successfully.
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go build: successful",
		projPath,
		[]string{"build", ".", "errors"},
	)
	// CI-Step2: check if the app tests run and pass as required.
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go test: successful",
		projPath,
		[]string{"test", "-v"},
	)

	// CI-Step3: check if the app code comforms with the golang formating rules.
	pipeline[2] = newExecutionStep(
		"go fmt",
		"gofmt",
		"Go fmt: successful",
		projPath,
		[]string{"-l", "."},
	)

	// CI-Step4: check if the app code comforms with the golang formating rules.
	pipeline[3] = newTimeoutStep(
		"git push",
		"git",
		"Git push: successful",
		projPath,
		[]string{"push", "origin", "main"},
		10*time.Second,
	)

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errChan <- err

				return
			}

			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errChan <- err

				return
			}
		}

		close(doneChan)
	}()

	for {
		select {
		case sig := <-sigChan:
			signal.Stop(sigChan)
			return fmt.Errorf("%s: Exiting: %w", sig, ErrSignal)
		case err := <-errChan:
			return err
		case <-doneChan:
			return nil
		}
	}
}
