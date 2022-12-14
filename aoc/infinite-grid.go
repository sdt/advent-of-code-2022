package aoc

import (
	"fmt"
)

type coord struct {
	x, y int
}

type InfiniteGrid[T any] struct {
	min, max coord
	defaultValue T
    cell map[coord]T
}

func (this InfiniteGrid[T]) Get(x, y int) T {
	if value, found := this.cell[coord{ x, y }]; found {
		return value
	}
	return this.defaultValue
}

func (this *InfiniteGrid[T]) Set(x, y int, value T) *InfiniteGrid[T] {
	if len(this.cell) == 0 {
		this.min.x = x
		this.min.y = y
		this.max.x = x
		this.max.y = y
	} else if x < this.min.x {
		this.min.x = x
	} else if x > this.max.x {
		this.max.x = x
	} else if y < this.min.y {
		this.min.y = y
	} else if y > this.max.y {
		this.max.y = y
	}

	this.cell[coord{x, y}] = value
	return this
}

func (this InfiniteGrid[T]) OnGrid(x, y int) bool {
	return x >= this.min.x && x <= this.max.x &&
	       y >= this.min.y && y <= this.max.y
}

func (this InfiniteGrid[T]) Print(format string, border int) {
	for y := this.min.y - border; y <= this.max.y + border; y++ {
		for x := this.min.x - border; x <= this.max.x + border; x++ {
			fmt.Printf(format, this.Get(x, y))
		}
		fmt.Println()
	}
	fmt.Println()
}

func NewInfiniteGrid[T any](defaultValue T) *InfiniteGrid[T] {
	return &InfiniteGrid[T]{
		defaultValue: defaultValue,
		cell: make(map[coord]T),
	}
}
