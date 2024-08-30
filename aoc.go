package goaoc

// Challenge represents the function signature expected for both parts of a given challenge.
// Each Challenge function receives a string input (raw challenge data) and returns an int result.
type Challenge func(string) int

// Part is an enumeration representing which part of the Advent of Code challenge to execute.
// Valid values are 1 and 2, corresponding to the problem statement's divisions.
type Part int

// NewPart constructs a Part from an integer. Returns an error if the part number is not valid (not 1 or 2).
//
// Example:
//
//	part, err := NewPart(2)
//	if err != nil {
//	    log.Fatal(err) // 'err' will contain 'invalid part' message if not 1 or 2
//	}
func NewPart(p int) (Part, error) {
	if p != 1 && p != 2 {
		return Part(0), InvalidPartError{Part: p}
	}

	return Part(p), nil
}
