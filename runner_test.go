package goaoc_test

import (
	"errors"
	"testing"

	"github.com/hvpaiva/goaoc"
	"github.com/hvpaiva/goaoc/mock"
)

func TestRunWithInvalidParts(t *testing.T) {
	testCases := []struct {
		name      string
		part      string
		expectErr string
	}{
		{"PartNotSpecified", "0", "invalid part: 0. The valid parts are (1/2)"},
		{"WrongPartDefined", "3", "invalid part: 3. The valid parts are (1/2)"},
		{"WrongPartTypeString", "ss", "invalid part type. The part type allowed is int"},
		{"WrongPartTypeEmpty", "", "invalid part type. The part type allowed is int"},
		{"WrongPartTypeStillString", "true", "invalid part type. The part type allowed is int"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mok := mock.NewManager(tc.part, nil, nil)
			err := goaoc.Run("input", mockPartOne, mockPartTwo, goaoc.WithManager(&mok))

			if err == nil || err.Error() != tc.expectErr {
				t.Fatalf("Expected error '%s', but got: %v", tc.expectErr, err)
			}
		})
	}
}

func TestRunWithErrors(t *testing.T) {
	testCases := []struct {
		name      string
		part      string
		selectErr error
		outputErr error
		expectErr string
	}{
		{"SelectingPartError", "2", errors.New("error when calling Read"), nil, "error when calling Read"},
		{"OutputError", "1", nil, errors.New("output failed"), "output failed"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mok := mock.NewManager(tc.part, tc.selectErr, tc.outputErr)
			err := goaoc.Run("input", mockPartOne, mockPartTwo, goaoc.WithManager(&mok))

			if err == nil || err.Error() != tc.expectErr {
				t.Fatalf("Expected error '%s', but got: %v", tc.expectErr, err)
			}
		})
	}
}

func TestRunWithValidPart(t *testing.T) {
	testCases := []struct {
		name           string
		part           string
		expectedOutput string
		copiedValue    string
	}{
		{"PartOne", "1", "The challenge result is 42\n", "42"},
		{"PartTwo", "2", "The challenge result is 24\n", "24"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mok := mock.NewManager(tc.part, nil, nil)
			err := goaoc.Run("input", mockPartOne, mockPartTwo, goaoc.WithManager(&mok))

			if err != nil {
				t.Fatalf("Unexpected error when part is valid: %v", err)
			}

			output := mok.GetStdout()
			expectedOutput := tc.expectedOutput
			if output != expectedOutput {
				t.Errorf("Expected output '%s', but got '%s'", expectedOutput, output)
			}
		})
	}
}

func TestRunWithDefaultManager(t *testing.T) {
	testCases := []struct {
		name string
		part int
	}{
		{"PartOne", 1},
		{"PartTwo", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := goaoc.Run("input", mockPartOne, mockPartTwo, goaoc.WithPart(tc.part))

			if err != nil {
				t.Fatalf("Unexpected error when part is valid: %v", err)
			}
		})
	}
}

func mockPartOne(_ string) int {
	return 42
}

func mockPartTwo(_ string) int {
	return 24
}
