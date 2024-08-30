package goaoc

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tiagomelo/go-clipboard/clipboard"
)

// Env struct embodies the input/output streams and command-line arguments used by IO managers.
// It provides flexibility in I/O handling by abstracting standard input and output mechanisms.
type Env struct {
	// Stdin is the input stream from which data can be read.
	// It typically corresponds to os.Stdin but can be overridden for testing or alternative input methods.
	Stdin io.Reader

	// Stdout is the output stream to which data can be written.
	// By default, it's set to os.Stdout but can be redirected to capture output programmatically or in tests.
	Stdout io.Writer

	// Args holds command-line arguments, minus the program name.
	// This slice allows the passing and manipulation of additional parameters through the command line.
	Args []string
}

var defaultConsoleEnv = Env{
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
	Args:   os.Args[1:],
}

// DefaultConsoleManager manages I/O via the default console, implementing IOManager.
type DefaultConsoleManager struct {
	Env Env
}

// NewConsoleManager initializes a new DefaultConsoleManager with standard console streams.
func NewConsoleManager() DefaultConsoleManager {
	return DefaultConsoleManager{
		Env: defaultConsoleEnv,
	}
}

// Read derives arguments like 'part' from various sources (flags, environment, or stdin).
// It returns errors if flag parsing fails or stdin input cannot be retrieved.
func (m DefaultConsoleManager) Read(arg string) (part string, err error) {
	if arg != "part" {
		return "", nil
	}

	checks := []func() (string, error){
		func() (string, error) { return getPartInFlag(m.Env) },
		getPartInEnv,
		func() (string, error) { return getPartInStdin(m.Env) },
	}

	for _, check := range checks {
		part, err = check()
		if err != nil {
			return "", err
		}

		if part != "" {
			return part, nil
		}
	}

	return part, IOReadError{Err: ErrMissingPart}
}

// Write outputs the result to console and optionally copies to clipboard if not disabled by GOAOC_DISABLE_COPY_CLIPBOARD.
// Errors can arise from console output failures or clipboard command errors.
func (m DefaultConsoleManager) Write(result string) error {
	if _, err := fmt.Fprintf(m.Env.Stdout, "The challenge result is %s\n", result); err != nil {
		return IOWriteError{Err: err}
	}

	toClipboard(result, m.Env.Stdout)

	return nil
}

// getPartInFlag attempts to parse the 'part' option from command-line flags.
// It supports standard flags only and returns errors if parsing fails.
func getPartInFlag(env Env) (part string, err error) {
	fs := flag.NewFlagSet("goaoc", flag.ContinueOnError)
	fs.SetOutput(env.Stdout)

	fs.Usage = func() {
		_, err = fmt.Fprintf(fs.Output(), "Usage: %s [options]\n", fs.Name())

		fs.PrintDefaults()
	}

	fs.StringVar(&part, "part", "", "Part of the challenge, valid values are (1/2)")

	if err = fs.Parse(env.Args); err != nil {
		return "", IOReadError{Err: err}
	}

	return part, nil
}

// getPartInEnv retrieves the 'part' from environment variables returned as a simple string.
func getPartInEnv() (string, error) {
	part := os.Getenv("GOAOC_CHALLENGE_PART")

	return part, nil
}

// getPartInStdin queries stdin to get which part the user wishes to run. Useful in interactive console mode.
// Returns errors for invalid or empty inputs.
func getPartInStdin(env Env) (string, error) {
	var part string

	_, err := fmt.Fprintln(env.Stdout, "Which part do you want to run? (1/2)")
	if err != nil {
		return "", err
	}

	_, err = fmt.Fscanln(env.Stdin, &part)
	if err != nil && errors.Is(err, io.EOF) {
		return "", IOReadError{Err: ErrMissingPart}
	}

	return part, nil
}

// toClipboard tries to copy the given value to the system clipboard. Skips copying if the environment is set to not copy.
// Errors while executing the clipboard command are printed but do not stop the program.
func toClipboard(value string, stdout io.Writer) {
	envVar := os.Getenv("GOAOC_DISABLE_COPY_CLIPBOARD")
	if envVar == "true" {
		return
	}

	c := clipboard.New()
	if err := c.CopyText(value); err != nil {
		_, _ = fmt.Fprintf(stdout, "Error copying to clipboard: %s\n", err)

		return
	}

	_, _ = fmt.Fprintf(stdout, "Copied to clipboard: %s\n", value)
}
