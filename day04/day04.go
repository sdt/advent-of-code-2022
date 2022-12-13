package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

type Assignment struct {
	first, last int
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	overlaps := 0
	for _, line := range lines {
		lhs, rhs := parseLine(line)
		if hasFullOverlap(lhs, rhs) {
			overlaps++
		}
	}
	return overlaps
}

func part2(lines []string) int {
	overlaps := 0
	for _, line := range lines {
		lhs, rhs := parseLine(line)
		if hasPartialOverlap(lhs, rhs) {
			overlaps++
		}
	}
	return overlaps
}

func parseLine(line string) (Assignment, Assignment) {
	assignments := strings.Split(line, ",")
	return parseAssignment(assignments[0]), parseAssignment(assignments[1])
}

func parseAssignment(line string) Assignment {
	values := aoc.ParseInts(strings.Split(line, "-"))

	return Assignment{values[0], values[1]}
}

func hasFullOverlap(lhs, rhs Assignment) bool {
	return ((lhs.first <= rhs.first) && (lhs.last >= rhs.last)) || ((rhs.first <= lhs.first) && (rhs.last >= lhs.last))
}

func hasPartialOverlap(lhs, rhs Assignment) bool {
	return !isDistinct(lhs, rhs)
}

func isDistinct(lhs, rhs Assignment) bool {
	return (lhs.last < rhs.first) || (rhs.last < lhs.first)
}
