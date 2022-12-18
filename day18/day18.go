package main

import (
	"advent-of-code/aoc"
	"fmt"
	"math"
	"strings"
)

type Vec3 [3]int

type BBox struct {
	min Vec3
	max Vec3
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	cubeMap := make(map[Vec3]bool)
	for _, line := range lines {
		cube := parseCube(line)
		cubeMap[cube] = true
	}

	fmt.Println(part1(cubeMap))
	fmt.Println(part2(cubeMap))
}

func part1(cubeMap map[Vec3]bool) int {
	surfaceArea := 0
	for pos := range cubeMap {
		for axis := range pos {
			before := pos
			before[axis]--
			if !cubeMap[before] {
				surfaceArea++
			}

			after := pos
			after[axis]++
			if !cubeMap[after] {
				surfaceArea++
			}
		}
	}
	return surfaceArea
}

func part2(cubeMap map[Vec3]bool) int {
	volume := makeVolume(cubeMap)
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

func makeVolume(cubes map[Vec3]bool) BBox {
	bbox := newBBox()
	for cube := range cubes {
		bbox = bbox.extend(cube)
	}

	for i := range bbox.min {
		bbox.min[i]--
		bbox.max[i]++
	}

	return bbox
}

func parseCube(line string) Vec3 {
	coords := aoc.ParseInts(strings.Split(line, ","))
	var p Vec3
	copy(p[:], coords)
	return p
}

func newBBox() BBox {
	min := [3]int{ math.MaxInt, math.MaxInt, math.MaxInt }
	max := [3]int{ math.MinInt, math.MinInt, math.MinInt }
	return BBox{ min, max }
}

func (this BBox) extend(p Vec3) BBox {
	var bbox BBox
	for axis, value := range p {
		bbox.min[axis] = min(this.min[axis], value)
		bbox.max[axis] = max(this.max[axis], value)
	}
	return bbox
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
