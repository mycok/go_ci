package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	testcases := []struct {
		name     string
		projPath string
		expOut   string
		expErr   error
		mockCmd  func(ctx context.Context, name string, args...string) *exec.Cmd
	}{
		// {
		// 	name:     "successful pipelines",
		// 	projPath: "./testdata/tool/",
		// 	expOut:   "Go build: successful\n"+
		// 				"Go test: successful\n"+
		// 				"Go fmt: successful\n"+
		// 				"Git push: successful\n",
		// 	expErr:   nil,
		// 	mockCmd: nil,
		// },
		{
			name:     "successful mock pipelines",
			projPath: "./testdata/tool/",
			expOut:   "Go build: successful\n"+
						"Go test: successful\n"+
						"Go fmt: successful\n"+
						"Git push: successful\n",
			expErr:   nil,
			mockCmd: mockCmdContext,
		},
		{
			name:     "failed build pipeline",
			projPath: "./testdata/toolErr",
			expOut:   "",
			expErr:   &stepErr{step: "go build"},
		},
		{
			name:     "failed fmt pipeline",
			projPath: "./testdata/toolFmtErr",
			expOut:   "",
			expErr:   &stepErr{step: "go fmt"},
		},
		{
			name:     "failed git pipeline with timeout on mocked cmd ",
			projPath: "./testdata/tool",
			expOut:   "",
			expErr:   context.DeadlineExceeded,
			mockCmd: mockCmdTimeout,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

			err := run(tc.projPath, &out)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q, but got 'nil' instead.", tc.expErr)

					return
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q, but got: %q instead.", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %q", err)

				return
			}

			if tc.expOut != out.String() {
				t.Errorf("Expected output msg to be: %q, but got %q instead", tc.expOut, out.String())
			}
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(11*time.Second)
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}

	os.Exit(1)
}

func mockCmdContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, name)
	cs = append(cs, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

	return cmd
}

func mockCmdTimeout(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, name, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")

	return cmd
}

