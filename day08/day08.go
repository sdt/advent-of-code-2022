package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Grid[T any] struct {
	w, h int
	cell []T
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	treeSize := ParseTrees(lines)
	visible := NewGrid[bool](treeSize.w, treeSize.h)

	numVisible := 0
	for y := 0; y < treeSize.h; y++ {
		numVisible += markVisible(treeSize, &visible, 0, y, 1, 0)
		numVisible += markVisible(treeSize, &visible, treeSize.w-1, y, -1, 0)
	}
	for x := 0; x < treeSize.w; x++ {
		numVisible += markVisible(treeSize, &visible, x, 0, 0, 1)
		numVisible += markVisible(treeSize, &visible, x, treeSize.h-1, 0, -1)
	}

	return numVisible
}

func part2(lines []string) int {
	treeSize := ParseTrees(lines)

	best := -1
	for y := 0; y < treeSize.h; y++ {
		for x := 0; x < treeSize.w; x++ {
			score := computeScenicScore(treeSize, x, y)
			if score > best {
				best = score
			}
		}
	}
	return best
}

func computeScenicScore(treeSize Grid[int], x, y int) int {
	u := computeSingleScore(treeSize, x, y, 0, -1)
	l := computeSingleScore(treeSize, x, y, -1, 0)
	r := computeSingleScore(treeSize, x, y, 1, 0)
	d := computeSingleScore(treeSize, x, y, 0, 1)
	score := u * l * r * d
	//fmt.Printf("%d,%d(%d) %d * %d * %d * %d = %d\n", x, y, treeSize.Get(x, y), u, l, r, d, score)
	return score
}

func computeSingleScore(treeSize Grid[int], x, y, dx, dy int) int {
	limit := treeSize.Get(x, y)
	count := 0
	for {
		x += dx
		y += dy
		if !treeSize.OnGrid(x, y) {
			return count
		}
		h := treeSize.Get(x, y)
		if h >= limit {
			return count + 1
		}
		count++
	}
	return count
}

func markVisible(treeSize Grid[int], visible *Grid[bool], x, y, dx, dy int) int {
	max := -1
	count := 0
	for treeSize.OnGrid(x, y) {
		h := treeSize.Get(x, y)
		if h <= max {
			//fmt.Printf("%d,%d is not visible (%d -> %d)\n", x, y, max, h)
		} else if visible.Get(x, y) {
			//fmt.Printf("%d,%d is already visible (%d -> %d)\n", x, y, max, h)
			max = h
		} else {
			//fmt.Printf("%d,%d is visible (%d -> %d)\n", x, y, max, h)
			visible.Set(x, y, true)
			max = h
			count++
		}
		x += dx
		y += dy
	}
	return count
}

func ParseTrees(lines []string) Grid[int] {
	height := NewGrid[int](len(lines[0]), len(lines))

	for y, line := range lines {
		for x, c := range line {
			height.Set(x, y, int(c-'0'))
		}
	}

	return height
}

func NewGrid[T any](w, h int) Grid[T] {
	return Grid[T]{
		w:    w,
		h:    h,
		cell: make([]T, w*h),
	}
}

func (this *Grid[T]) Set(x, y int, value T) {
	this.cell[this.offset(x, y)] = value
}

func (this *Grid[T]) Get(x, y int) T {
	return this.cell[this.offset(x, y)]
}

func (this Grid[T]) OnGrid(x, y int) bool {
	return (x >= 0) && (y >= 0) && (x < this.w) && (y < this.h)
}

func (this Grid[T]) offset(x, y int) int {
	return y*this.w + x
}
