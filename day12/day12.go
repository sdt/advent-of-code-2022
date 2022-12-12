package main

import (
    "advent-of-code/aoc"
    "fmt"
)

type Vec2 struct {
    x, y int
}

type Grid[T any] struct {
    size Vec2
    cell []T
}

type Queue[T any] struct {
    queue []T
}

type Solution struct {
    pos Vec2
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
    grid := NewGrid[int](len(lines[0]) + 2, len(lines) + 2, 100)
    seen := NewGrid[bool](grid.size.x, grid.size.y, false)

    var end Vec2

    agenda := NewQueue[Solution]()

    directions := []Vec2{ Vec2{ 1, 0 }, Vec2{ 0, 1 }, Vec2{ -1, 0 }, Vec2{ 0, -1 } }
    for y, line := range lines {
        for x, char := range line {
            pos := Vec2{ x + 1, y + 1 }

            if char == 'S' {
                char = 'a'
                agenda.Enqueue(NewSolution(pos, 0))
                seen.Set(pos, true)
            } else if char == 'E' {
                char = 'z'
                end = pos
            }

            grid.Set(pos, int(char - 'a'))
        }
    }

    for !agenda.IsEmpty() {
        attempt := agenda.Dequeue()
        if attempt.pos == end {
            return attempt.distance
        }


        from := grid.Get(attempt.pos)
        for _, direction := range directions {
            next := attempt.pos.Add(direction)
            if seen.Get(next) {
                continue;
            }
            to := grid.Get(next)

            if to - from <= 1 {
                agenda.Enqueue(NewSolution(next, attempt.distance+1))
                seen.Set(next, true)
            }
        }
    }

    return -1
}

func part2(lines []string) int {
    // Leave a border around the edge as sentinel values
    grid := NewGrid[int](len(lines[0]) + 2, len(lines) + 2, -100)
    seen := NewGrid[bool](grid.size.x, grid.size.y, false)

    agenda := NewQueue[Solution]()

    directions := []Vec2{ Vec2{ 1, 0 }, Vec2{ 0, 1 }, Vec2{ -1, 0 }, Vec2{ 0, -1 } }
    for y, line := range lines {
        for x, char := range line {
            pos := Vec2{ x + 1, y + 1 }

            if char == 'S' {
                char = 'a'
            } else if char == 'E' {
                char = 'z'
                agenda.Enqueue(NewSolution(pos, 0))
                seen.Set(pos, true)
            }

            grid.Set(pos, int(char - 'a'))
        }
    }

    for !agenda.IsEmpty() {
        attempt := agenda.Dequeue()

        from := grid.Get(attempt.pos)
        if from == 0 {
            return attempt.distance
        }

        for _, direction := range directions {
            next := attempt.pos.Add(direction)
            if seen.Get(next) {
                continue;
            }
            to := grid.Get(next)

            if from - to <= 1 {
                agenda.Enqueue(NewSolution(next, attempt.distance+1))
                seen.Set(next, true)
            }
        }
    }
    return -1
}

func NewGrid[T any](w, h int, defaultValue T) *Grid[T] {
    grid := Grid[T]{
        size: Vec2{w, h},
        cell: make([]T, w * h),
    }

    for i := range grid.cell {
        grid.cell[i] = defaultValue
    }

    return &grid
}

func (this *Grid[T]) Set(pos Vec2, value T) *Grid[T] {
    this.cell[this.offset(pos)] = value
    return this
}

func (this *Grid[T]) Get(pos Vec2) T {
    return this.cell[this.offset(pos)]
}

func (this Grid[T]) offset(pos Vec2) int {
    return this.size.x * pos.y + pos.x
}

func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{ queue: make([]T, 0) }
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
    return Solution{ pos: pos, distance: distance }
}

func (this Vec2) Add(that Vec2) Vec2 {
    return Vec2{ x: this.x + that.x, y: this.y + that.y }
}
