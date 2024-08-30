// Copyright (c) 2024 Highlander Paiva. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package goaoc

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func mockEnv(args []string, input string, output io.Writer) Env {
	return Env{
		Stdin:  bytes.NewBufferString(input),
		Stdout: output,
		Args:   args,
	}
}

type failingWriter struct{}

func (f *failingWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("write failed")
}

func TestRead(t *testing.T) {
	testCases := []struct {
		name      string
		env       Env
		expect    string
		expectErr string
	}{
		{"PartFromFlag", mockEnv([]string{"-part=1"}, "", new(bytes.Buffer)), "1", ""},
		{"PartFromEnv", mockEnv([]string{}, "", new(bytes.Buffer)), "2", ""},
		{"PartFromStdin", mockEnv([]string{}, "1\n", new(bytes.Buffer)), "1", ""},
		{"PartFromStdinFailStdout", mockEnv([]string{}, "1\n", &failingWriter{}), "1", "write failed"},
		{"PartFromStdinFailEmpty", mockEnv([]string{}, "", new(bytes.Buffer)), "", "failed to read input: no part specified, please provide a valid part"},
		{"FlagProvidedButNotDefined", mockEnv([]string{"--test"}, "0", new(bytes.Buffer)), "", "failed to read input: flag provided but not defined: -test"},
		{"FlagProvidedButNotDefinedFailedStdout", mockEnv([]string{"--test"}, "0", &failingWriter{}), "", "failed to read input: flag provided but not defined: -test"},
		{"EmptyRead", mockEnv([]string{}, "", new(bytes.Buffer)), "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager := DefaultConsoleManager{Env: tc.env}

			if tc.name == "PartFromEnv" {
				_ = os.Setenv("GOAOC_CHALLENGE_PART", "2")
				defer func() {
					err := os.Unsetenv("GOAOC_CHALLENGE_PART")
					if err != nil {
						t.Fatalf("Unexpected error while unsetting environment variable: %v", err)
					}
				}()
			}

			part, err := manager.Read("part")
			if tc.name == "EmptyRead" {
				part, err = manager.Read("")

			}
			if tc.expectErr != "" {
				if err == nil || err.Error() != tc.expectErr {
					t.Fatalf("Expected error '%s', but got: %v", tc.expectErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if part != tc.expect {
					t.Errorf("Expected part %s, but got %s", tc.expect, part)
				}
			}
		})
	}
}

func TestToClipboard(t *testing.T) {
	env := mockEnv([]string{}, "", new(bytes.Buffer))
	manager := DefaultConsoleManager{Env: env}

	testCases := []struct {
		name   string
		output string
	}{
		{"Working", "Copied to clipboard: test value"},
		{"Deactivated", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Setenv("GOAOC_DISABLE_COPY_CLIPBOARD", "false")
			if tc.name == "Deactivated" {
				_ = os.Setenv("GOAOC_DISABLE_COPY_CLIPBOARD", "true")
			}

			toClipboard("test value", env.Stdout)

			output := manager.Env.Stdout.(*bytes.Buffer).String()
			if !strings.Contains(output, tc.output) {
				t.Errorf("Expected clipboard message, but got: %s", output)
			}
		})
	}
}

func TestOutput(t *testing.T) {
	env := mockEnv([]string{}, "", new(bytes.Buffer))
	manager := DefaultConsoleManager{Env: env}
	_ = os.Setenv("GOAOC_DISABLE_COPY_CLIPBOARD", "false")

	err := manager.Write("42")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := env.Stdout.(*bytes.Buffer).String()
	expectedOutput := "The challenge result is 42\nCopied to clipboard: 42\n"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', but got '%s'", expectedOutput, output)
	}
}

func TestSelectPartErrors(t *testing.T) {
	_ = os.Unsetenv("GOAOC_CHALLENGE_PART")

	mockEnv := mockEnv([]string{}, "", new(bytes.Buffer))

	manager := DefaultConsoleManager{Env: mockEnv}

	_, err := manager.Read("part")
	if err == nil || err.Error() != "failed to read input: no part specified, please provide a valid part" {
		t.Fatalf("Expected 'failed to read input: no part specified, please provide a valid part' error, but got: %v", err)
	}
}

func TestOutputWriterFails(t *testing.T) {
	manager := DefaultConsoleManager{
		Env: Env{
			Stdout: &failingWriter{},
		},
	}

	err := manager.Write("42")
	if err == nil || err.Error() != "failed to write input: write failed" {
		t.Fatalf("Expected 'failed to write input: write failed' error, but got: %v", err)
	}
}

func TestNewConsoleManager(t *testing.T) {
	manager := NewConsoleManager()

	if !reflect.DeepEqual(manager.Env, defaultConsoleEnv) {
		t.Errorf("expected Stdin to be %v, but got %v", os.Stdin, manager.Env.Stdin)
	}
}
