package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		total += getPriorities(line)
	}
	return total
}

func part2(lines []string) int {
	total := 0
	for i := 0; i < len(lines); i += 3 {
		total += getCommon(lines[i : i+3])
	}
	return total
}

func getPriorities(rucksack string) int {
	seen := make([]int, 52)

	half := len(rucksack) / 2
	for i := 0; i < half; i++ {
		item := rucksack[i]
		value := getValue(item)
		seen[value] = 1
	}
	for i := half; i < len(rucksack); i++ {
		item := rucksack[i]
		value := getValue(item)
		if seen[value] > 0 {
			return value + 1
		}
	}
	return 0
}

func getCommon(lines []string) int {
	seen := make([]int, 52)

	all := 0
	for i, line := range lines {
		bit := 1 << i
		all |= bit

		for _, item := range line {
			value := getValue(byte(item))
			seen[value] |= bit
		}
	}

	for i, mask := range seen {
		if mask == all {
			return i + 1
		}
	}
	return 0
}

func getValue(item byte) int {
	if item >= 'a' && item <= 'z' {
		return int(item - 'a')
	}
	return int(item - 'A' + 26)
}
