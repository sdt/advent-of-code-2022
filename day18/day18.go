package main

import (
	"advent-of-code/aoc"
	"fmt"
	"sort"
	"strings"
)

type Side uint8

// RH coords.
const (
	East     = 1 << iota // +x
	West     = 1 << iota // -x
	North    = 1 << iota // +y
	South    = 1 << iota // -y
	Up       = 1 << iota // +Z
	Down     = 1 << iota // -Z

	XMask    = East | West
	YMask    = North | South
	ZMask    = Up | Down

	All		 = XMask | YMask | ZMask
)

type Cube struct {
	coords	  [3]int
	openSides Side
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	cubes := make([]*Cube, len(lines))

	for i, line := range lines {
		cubes[i] = parseCube(line)
	}

	for axis := 0; axis < 3; axis++ {
		matchCubes(cubes, axis)
	}

	area := 0
	for _, cube := range(cubes) {
		area += cube.SurfaceArea()
	}

	return area
}

func matchCubes(cubes []*Cube, axis int) {
	sort.Slice(cubes, func(a, b int) bool {
		return cubes[a].coords[axis] < cubes[b].coords[axis]
	})

	for i, a := range cubes[:len(cubes)-1] {
		for _, b := range cubes[i+1:] {
			if isAdjacent(a, b, axis) {
				aSide := Side(1 << (axis * 2))
				bSide := Side(aSide << 1)

				a.openSides &= ^aSide
				b.openSides &= ^bSide
			}
		}
	}
}

func isAdjacent(a, b *Cube, axis int) bool {
	for c := 0; c < 3; c++ {
		if c == axis {
			if a.coords[c] != b.coords[c] - 1 {
				return false
			}
		} else {
			if a.coords[c] != b.coords[c] {
				return false
			}
		}
	}
	return true
}

func parseCube(line string) *Cube {
	coords := aoc.ParseInts(strings.Split(line, ","))
	cube := Cube{ openSides: All }
	copy(cube.coords[:], coords)
	return &cube
}

func (this Cube) SurfaceArea() int {
	area := 0
	for i := 0; i < 6; i++ {
		if this.openSides & (1 << i) != 0 {
			area++
		}
	}
	return area
}
