package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Vec2 struct {
	x, y int
}

type Queue[T any] struct {
	queue []T
}

type Solution struct {
	pos      Vec2
	distance int
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	// Leave a border around the edge as sentinel values
	grid := aoc.NewGrid[int](len(lines[0])+2, len(lines)+2, 100)
	seen := aoc.NewGrid[bool](grid.Width(), grid.Height(), false)

	var end Vec2

	agenda := NewQueue[Solution]()

	directions := []Vec2{Vec2{1, 0}, Vec2{0, 1}, Vec2{-1, 0}, Vec2{0, -1}}
	for y, line := range lines {
		for x, char := range line {
			pos := Vec2{x + 1, y + 1}

			if char == 'S' {
				char = 'a'
				agenda.Enqueue(NewSolution(pos, 0))
				seen.Set(pos.x, pos.y, true)
			} else if char == 'E' {
				char = 'z'
				end = pos
			}

			grid.Set(pos.x, pos.y, int(char-'a'))
		}
	}

	for !agenda.IsEmpty() {
		attempt := agenda.Dequeue()
		if attempt.pos == end {
			return attempt.distance
		}

		from := grid.Get(attempt.pos.x, attempt.pos.y)
		for _, direction := range directions {
			next := attempt.pos.Add(direction)
			if seen.Get(next.x, next.y) {
				continue
			}
			to := grid.Get(next.x, next.y)

			if to-from <= 1 {
				agenda.Enqueue(NewSolution(next, attempt.distance+1))
				seen.Set(next.x, next.y, true)
			}
		}
	}

	return -1
}

func part2(lines []string) int {
	// Leave a border around the edge as sentinel values
	grid := aoc.NewGrid[int](len(lines[0])+2, len(lines)+2, -100)
	seen := aoc.NewGrid[bool](grid.Width(), grid.Height(), false)

	agenda := NewQueue[Solution]()

	directions := []Vec2{Vec2{1, 0}, Vec2{0, 1}, Vec2{-1, 0}, Vec2{0, -1}}
	for y, line := range lines {
		for x, char := range line {
			pos := Vec2{x + 1, y + 1}

			if char == 'S' {
				char = 'a'
			} else if char == 'E' {
				char = 'z'
				agenda.Enqueue(NewSolution(pos, 0))
				seen.Set(pos.x, pos.y, true)
			}

			grid.Set(pos.x, pos.y, int(char-'a'))
		}
	}

	for !agenda.IsEmpty() {
		attempt := agenda.Dequeue()

		from := grid.Get(attempt.pos.x, attempt.pos.y)
		if from == 0 {
			return attempt.distance
		}

		for _, direction := range directions {
			next := attempt.pos.Add(direction)
			if seen.Get(next.x, next.y) {
				continue
			}
			to := grid.Get(next.x, next.y)

			if from-to <= 1 {
				agenda.Enqueue(NewSolution(next, attempt.distance+1))
				seen.Set(next.x, next.y, true)
			}
		}
	}
	return -1
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{queue: make([]T, 0)}
}

func (this *Queue[T]) Enqueue(value T) *Queue[T] {
	this.queue = append(this.queue, value)
	return this
}

func (this *Queue[T]) Dequeue() T {
	ret := this.queue[0]
	this.queue = this.queue[1:]
	return ret
}

func (this Queue[T]) IsEmpty() bool {
	return len(this.queue) == 0
}

func NewSolution(pos Vec2, distance int) Solution {
	return Solution{pos: pos, distance: distance}
}

func (this Vec2) Add(that Vec2) Vec2 {
	return Vec2{x: this.x + that.x, y: this.y + that.y}
}
