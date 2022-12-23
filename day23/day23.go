package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Vec2 struct {
	x, y int
}

var (
	N = Vec2{0, -1}
	E = Vec2{1, 0}
	S = Vec2{0, 1}
	W = Vec2{-1, 0}

	NE = N.Add(E)
	NW = N.Add(W)
	SE = S.Add(E)
	SW = S.Add(W)
)

var directions = [][]Vec2{
	{N, NE, NW},
	{S, SE, SW},
	{W, NW, SW},
	{E, NE, SE},
}

type ElfMap map[Vec2]bool

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	current := parseInput(lines)
	//current.Print()

	for round := 0; round < 10; round++ {
		next, changed := step(current, round%len(directions))
		//fmt.Println("After round", round+1)
		if !changed {
			//fmt.Println("No more changes")
			break
		}
		//next.Print()
		current = next
	}

	min, max := current.BoundingBox()
	total := (max.x - min.x + 1) * (max.y - min.y + 1)

	return total - len(current)
}

func part2(lines []string) int {
	current := parseInput(lines)
	//current.Print()

	for round := 0; ; round++ {
		next, changed := step(current, round%len(directions))
		//fmt.Println("After round", round+1)
		if !changed {
			//fmt.Println("No more changes")
			return round + 1
		}
		//next.Print()
		current = next
	}
}

func step(current ElfMap, heading int) (ElfMap, bool) {
	count := make(map[Vec2]int)
	next := make(ElfMap, len(current))
	changed := false

	// First half
	for pos := range current {
		if hasNeighbour(current, pos) {
			nextPos := proposeMove(current, pos, heading)
			if nextPos == pos {
				//fmt.Printf("%v has no proposed move and isn't moving\n", pos)
				next[pos] = true // no move -> stay still
			} else {
				//fmt.Printf("%v proposed moving to %v\n", pos, nextPos)
				count[nextPos] = count[nextPos] + 1
			}
		} else {
			//fmt.Printf("%v has no neighbours and isn't moving\n", pos)
			next[pos] = true // elf with no neighours doesn't move
		}
	}

	// Second half
	for pos := range current {
		if next[pos] {
			// this elf isn't moving
			continue
		}
		nextPos := proposeMove(current, pos, heading)
		if count[nextPos] == 1 {
			//fmt.Printf("%v fulfils proposed move to %v\n", pos, nextPos)
			next[nextPos] = true // move
			changed = true
		} else {
			//fmt.Printf("%v abandons proposed move\n", pos)
			next[pos] = true // no move
		}
	}

	return next, changed
}

func hasNeighbour(elfMap ElfMap, pos Vec2) bool {
	var delta Vec2
	for delta.y = -1; delta.y <= 1; delta.y++ {
		for delta.x = -1; delta.x <= 1; delta.x++ {
			n := pos.Add(delta)
			if n != pos && elfMap[n] {
				return true
			}
		}
	}
	return false
}

func proposeMove(elfMap ElfMap, pos Vec2, heading int) Vec2 {
	for i := 0; i < len(directions); i++ {
		dirs := directions[(i+heading)%len(directions)]
		if allClear(elfMap, pos, dirs) {
			return pos.Add(dirs[0])
		}
	}
	return pos
}

func allClear(elfMap ElfMap, pos Vec2, dirs []Vec2) bool {
	for _, dir := range dirs {
		if elfMap[pos.Add(dir)] {
			return false
		}
	}
	return true
}

func parseInput(lines []string) ElfMap {
	elfMap := make(ElfMap)

	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				p := Vec2{x, y}
				elfMap[p] = true
			}
		}
	}

	return elfMap
}

func (this Vec2) Add(that Vec2) Vec2 {
	return Vec2{this.x + that.x, this.y + that.y}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (this Vec2) Min(that Vec2) Vec2 {
	return Vec2{min(this.x, that.x), min(this.y, that.y)}
}

func (this Vec2) Max(that Vec2) Vec2 {
	return Vec2{max(this.x, that.x), max(this.y, that.y)}
}

func (this ElfMap) BoundingBox() (Vec2, Vec2) {
	big := 1_000_000
	min := Vec2{+big, +big}
	max := Vec2{-big, -big}

	for pos := range this {
		min = min.Min(pos)
		max = max.Max(pos)
	}

	return min, max
}

func (this ElfMap) Print() {
	min, max := this.BoundingBox()

	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			if this[Vec2{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
