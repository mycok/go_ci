package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testcases := []struct{
		name string
		projPath string
		expOut string
		expErr error
	}{
		{
			name:    "successful build",
			projPath: "./testdata/tool",
			expOut:  "Go Build: Successful\n",
			expErr:  nil,
		},
		{
			name:    "failed build",
			projPath: "./testdata/toolErr",
			expOut:  "",
			expErr:  &stepErr{step: "go build"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			err := run(tc.projPath, &out)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q, but got 'nil' instead.", tc.expErr)

					return
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q, but got: %v instead.", tc.expErr, err)
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