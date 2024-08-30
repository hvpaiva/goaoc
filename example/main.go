package main

import (
	"log"

	"github.com/hvpaiva/goaoc"
)

func main() {
	err := goaoc.Run("input", partOne, partTwo)
	if err != nil {
		log.Fatalf("error running Go AoC: %v", err)
	}
}

func partOne(input string) int {
	return len(input)
}

func partTwo(input string) int {
	return len(input) * 2
}
