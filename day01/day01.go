package main

import (
	"advent-of-code/aoc"
	"fmt"
	"sort"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	largest := 0
	current := 0

	for _, line := range lines {
		if line == "" {
			current = 0
		} else {
			calories := aoc.ParseInt(line)
			current += calories
			if current > largest {
				largest = current
			}
		}
	}
	return largest
}

func part2(lines []string) int {
	totals := make([]int, 1)

	for _, line := range lines {
		if line == "" {
			totals = append(totals, 0)
		} else {
			calories := aoc.ParseInt(line)
			totals[len(totals)-1] += calories
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(totals)))
	return totals[0] + totals[1] + totals[2]
}
