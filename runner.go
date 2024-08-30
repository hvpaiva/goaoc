// Copyright (c) 2024 Highlander Paiva. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

// Package goaoc provides a framework to facilitate running Advent of Code challenges.
// This package encompasses utilities for handling input and outputs, selecting challenge parts,
// and executing them with configurable options.
//
// # Overview
//
// The goaoc package is designed to simplify the execution of Advent of Code challenges
// by abstracting I/O operations and enabling easy switching between parts 1 and 2
// of each challenge. The main entry point for executing a challenge is the Run function.
//
// # Basic Usage
//
// To run a challenge, you need to provide input data and two functions implementing
// the Challenge type, each corresponding to part 1 and part 2 of the challenge.
//
// Example:
//
//	err := Run("yourInputData", part1Func, part2Func, WithPart(1))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Additional RunOptions such as WithManager and WithPart allow customization of
// input/output management and challenge part selection, respectively.
package goaoc

import (
	"strconv"
)

// runOptions holds the configurations needed for running a challenge.
// It includes the IOManager for handling input/output and the challenge Part.
type runOptions struct {
	manager IOManager
	part    Part
}

// RunOption is a functional option type for configuring runOptions.
// It allows the user to customize aspects of the Run function.
type RunOption func(options *runOptions) error

// IOManager is an interface that abstracts the process of reading and writing data.
// It allows for different implementations to manage input and output according to varying needs, such as
// console-based, file-based, or even network-based I/O.
type IOManager interface {
	// Write writes the result string to an output destination.
	// Implementations must handle errors that occur during the write operation, such as IO errors.
	// Example:
	//   err := manager.Write("result data")
	//   if err != nil {
	//       log.Println("Failed to write result:", err)
	//   }
	Write(result string) error

	// Read retrieves a value based on the given argument string.
	// It's typically used to fetch configuration settings like which part of a challenge to run.
	// Errors may result from issues such as missing data or failed parse attempts.
	// Example:
	//   arg, err := manager.Read("part")
	//   if err != nil {
	//       log.Println("Failed to read argument:", err)
	//   }
	Read(arg string) (string, error)
}

// Run executes given Challenge functions partOne and partTwo, based on the input provided
// and optional configurations. It writes output via the configured IOManager.
//
// Example:
//
//	err := Run("123", func(input string) int { return len(input) }, func(input string) int { return len(input) * 2 }, WithPart(1))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// By default, output is written to the console, but you can change this by providing different IOManagers.
//
// Possible errors include option injection failures, I/O errors, and invalid part errors.
func Run(input string, partOne, partTwo Challenge, options ...RunOption) error {
	var opts runOptions
	if err := injectOptions(&opts, options...); err != nil {
		return err
	}

	result := executeChallenge(input, partOne, partTwo, opts.part)

	if err := opts.manager.Write(strconv.Itoa(result)); err != nil {
		return err
	}

	return nil
}

// WithManager creates a RunOption to set the custom IOManager.
// Use this to override the default console-based manager.
//
// Example:
//
//	manager := NewCustomManager()
//	err := Run(inputData, part1Func, part2Func, WithManager(manager))
func WithManager(manager IOManager) RunOption {
	return func(options *runOptions) error {
		options.manager = manager

		return nil
	}
}

// WithPart creates a RunOption to specify which part of the challenge to run (part 1 or 2).
// This is particularly useful when you want to determine the part dynamically.
//
// Example:
//
//	err := Run(inputData, part1Func, part2Func, WithPart(2))
func WithPart(part int) RunOption {
	return func(options *runOptions) error {
		options.part = Part(part)

		return nil
	}
}

// executeChallenge applies the appropriate Challenge function based on the selected part.
// It returns the result of the challenge execution.
func executeChallenge(input string, partOne, partTwo Challenge, part Part) (result int) {
	switch part {
	case 1:
		result = partOne(input)
	case 2:
		result = partTwo(input)
	default:
		// Though should never reach, it is good for future-proofing
		panic(ErrMissingPart)
	}

	return result
}

// injectOptions applies the functional options to configure runOptions.
// It defaults the IOManager to a console manager and resolves the challenge part from input if not set.
func injectOptions(opts *runOptions, options ...RunOption) error {
	for _, option := range options {
		_ = option(opts)
	}

	if opts.manager == nil {
		opts.manager = NewConsoleManager()
	}

	if opts.part == 0 {
		partStr, err := opts.manager.Read("part")
		if err != nil {
			return err
		}

		part, err := strconv.Atoi(partStr)
		if err != nil {
			return ErrInvalidPartType
		}

		opts.part, err = NewPart(part)
		if err != nil {
			return err
		}
	}

	return nil
}
