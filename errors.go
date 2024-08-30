package goaoc

import (
	"errors"
	"fmt"
)

// InvalidPartError indicates an error that occurs when an invalid part number
// is specified. Valid part numbers are 1 and 2.
type InvalidPartError struct {
	Part int
}

// Error implements the error interface for InvalidPartError.
// It returns a descriptive error message suitable for logging and debugging.
func (e InvalidPartError) Error() string {
	return fmt.Sprintf("invalid part: %d. The valid parts are (1/2)", e.Part)
}

// ErrInvalidPartType indicates an error that occurs when an invalid part type
// is specified. Valid part type is int.
var ErrInvalidPartType = errors.New("invalid part type. The part type allowed is int")

// ErrMissingPart indicates that no part was specified when it is required.
// This error typically occurs during input parsing when the part number
// is expected to be provided by some means (flag, input, etc.).
var ErrMissingPart = errors.New("no part specified, please provide a valid part")

// IOReadError indicates a failure during input operations, such as reading
// from a file or receiving input from the console. The underlying error
// can be retrieved for detailed inspection if necessary.
type IOReadError struct {
	Err error
}

// Error implements the error interface for IOReadError.
// It provides a message indicating an I/O read failure.
func (e IOReadError) Error() string {
	return fmt.Sprintf("failed to read input: %v", e.Err)
}

// Unwrap allows access to the underlying error, following Go 1.13's error unwrapper design.
func (e IOReadError) Unwrap() error {
	return e.Err
}

// IOWriteError indicates a failure during output operations, such as writing
// to a file or console. The underlying error
// can be retrieved for detailed inspection if necessary.
type IOWriteError struct {
	Err error
}

// Error implements the error interface for IOWriteError.
// It provides a message indicating an I/O write failure.
func (e IOWriteError) Error() string {
	return fmt.Sprintf("failed to write input: %v", e.Err)
}

// Unwrap allows access to the underlying error, following Go 1.13's error unwrapper design.
func (e IOWriteError) Unwrap() error {
	return e.Err
}
