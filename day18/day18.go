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

type Vec3 [3]int

type Cube struct {
	coords	  Vec3
	openSides Side
}

type BBox struct {
	min Vec3
	max Vec3
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	cubes := make([]*Cube, len(lines))
	for i, line := range lines {
		cubes[i] = parseCube(line)
	}

	fmt.Println(part1(cubes))
	fmt.Println(part2(cubes))
}

func part1(cubes []*Cube) int {
	return findSurfaceArea(cubes)
}

func part2(cubes []*Cube) int {
	cubeMap := make(map[Vec3]bool)
	for _, cube := range cubes {
		cubeMap[cube.coords] = true
	}

	volume := makeVolume(cubes)
	seen := make(map[Vec3]bool)
	surfaceArea := 0

	checkCell := func(pos Vec3) bool {
		if !volume.contains(pos) {
			return false
		}

		if seen[pos] {
			return false
		}

		if cubeMap[pos] {
			// Gone from empty space to inside cube. We've crossed a new surface
			// area face.
			surfaceArea++
			return false
		}

		// Still in unseen empty space
		return true
	}

	var floodFill func(Vec3)
	floodFill = func(pos Vec3) {
		seen[pos] = true // mark this one off

		for axis := range pos {
			before := pos
			before[axis]--
			if checkCell(before) {
				floodFill(before)
			}

			after := pos
			after[axis]++
			if checkCell(after) {
				floodFill(after)
			}
		}
	}
	floodFill(volume.min)

	return surfaceArea
}

func findSurfaceArea(cubes []*Cube) int {
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

func makeVolume(cubes []*Cube) BBox {
	bbox := newBBox(*cubes[0])
	for _, cube := range cubes[1:] {
		bbox = bbox.extend(*cube)
	}

	for i := range bbox.min {
		bbox.min[i]--
		bbox.max[i]++
	}

	return bbox
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

func newBBox(cube Cube) BBox {
	return BBox{ cube.coords, cube.coords }
}

func (this BBox) extend(cube Cube) BBox {
	var bbox BBox
	for axis, value := range cube.coords {
		bbox.min[axis] = min(this.min[axis], value)
		bbox.max[axis] = max(this.max[axis], value)
	}
	return bbox
}

func (this BBox) size(axis int) int {
	return this.max[axis] - this.min[axis] + 1
}

func (this BBox) contains(p Vec3) bool {
	for axis, value := range p {
		if value < this.min[axis] || value > this.max[axis] {
			return false
		}
	}
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
