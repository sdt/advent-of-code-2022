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
}

func part1(lines[] string) int {
    treeSize := ParseTrees(lines)
    visible := NewGrid[bool](treeSize.w, treeSize.h)

    numVisible := 0
    for y := 0; y < treeSize.h; y++ {
        numVisible += markVisible(treeSize, &visible, 0,            y, 1, 0)
        numVisible += markVisible(treeSize, &visible, treeSize.w-1, y, -1, 0)
    }
    for x := 0; x < treeSize.w; x++ {
        numVisible += markVisible(treeSize, &visible, x, 0,            0,  1)
        numVisible += markVisible(treeSize, &visible, x, treeSize.h-1, 0, -1)
    }

    return numVisible
}

func markVisible(treeSize Grid[int], visible *Grid[bool], x, y, dx, dy int) int {
    max := -1
    count := 0
    for treeSize.onGrid(x, y) {
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
            height.Set(x, y, int(c - '0'))
        }
    }

    return height
}

func NewGrid[T any](w, h int) Grid[T] {
    return Grid[T]{
        w: w,
        h: h,
        cell: make([]T, w * h),
    }
}

func (this *Grid[T]) Set(x, y int, value T) {
    this.cell[this.offset(x, y)] = value
}

func (this *Grid[T]) Get(x, y int) T {
    return this.cell[this.offset(x, y)]
}

func (this Grid[T]) onGrid(x, y int) bool {
    return (x >= 0) && (y >= 0) && (x < this.w) && (y < this.h)
}

func (this Grid[T]) offset(x, y int) int {
    return y * this.w + x
}
