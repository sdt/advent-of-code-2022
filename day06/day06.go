package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	filename := aoc.GetFilename()
	input := aoc.Slurp(filename)

	fmt.Println(part1(input))
	fmt.Println(part2(input))
}

func part1(input string) int {
	return findMarker(input, 4)
}

func part2(input string) int {
	return findMarker(input, 14)
}

func findMarker(input string, markerLen int) int {
	lastPos := len(input) - markerLen + 1

	for i := 0; i < lastPos; i++ {
		candidate := input[i : i+markerLen]
		if isMarker(candidate) {
			return i + markerLen
		}
	}

	return 0
}

func isMarker(candidate string) bool {
	mask := uint32(0)

	for _, char := range candidate {
		bit := uint32(1) << (char - 'a')

		if mask|bit == mask {
			return false
		}
		mask |= bit
	}
	return true
}
