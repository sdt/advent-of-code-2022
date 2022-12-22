package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Tile rune
const (
	Empty Tile = ' '
	Open  Tile = '.'
	Wall  Tile = '#'
)

type Span struct {
	offset, size int
}

type Map struct {
	tiles Grid[Tile]
	rows []Span
	cols []Span
}

type Vec2 struct {
	x, y int
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	moves := lines[len(lines)-1]

	lines = lines[:len(lines)-2]
	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	m := NewMap(width, height)
	for y, line := range lines {
		for x, char := range line {
			m.tiles.Set(x, y, Tile(char))
		}
	}

	// Build row spans
	for y := 0; y < height; y++ {
		x := 0
		for m.tiles.Get(x, y) == Empty {
			x++
		}
		m.rows[y].offset = x
		for m.tiles.Get(x, y) != Empty {
			x++
		}
		m.rows[y].size = x - m.rows[y].offset
	}

	// Build col spans
	for x := 0; x < width; x++ {
		y := 0
		for m.tiles.Get(x, y) == Empty {
			y++
		}
		m.cols[x].offset = y
		for m.tiles.Get(x, y) != Empty {
			y++
		}
		m.cols[x].size = y - m.cols[x].offset
	}

	move := [...]func(*Map, Vec2, int) Vec2{
		MoveRight,
		MoveDown,
		MoveLeft,
		MoveUp,
	}
	//dirname := [...]string{ "right", "down", "left", "up" }

	pos := Vec2{ m.rows[0].offset, 0 }
	dir := 0

	var distance, turn int

	for len(moves) > 0 {
		distance, moves = getMoveDistance(moves)
		newPos := move[dir](&m, pos, distance)
		//fmt.Printf("Move %d %s from %d,%d to %d,%d\n", distance, dirname[dir], pos.x, pos.y, newPos.x, newPos.y)
		pos = newPos

		if len(moves) == 0 {
			break
		}

		turn, moves = getTurn(moves)
		dir = (dir + len(move) + turn) % len(move)
	}

	//fmt.Printf("Final row=%d col=%d dir=%d\n", pos.y + 1, pos.x + 1, dir)

	return 1000 * (pos.y + 1) + 4 * (pos.x + 1) + dir
}

func MoveUp(m *Map, pos Vec2, count int) Vec2 {
	// count %= m.rows[pos.y].size
	col := m.cols[pos.x]
	for i := 0; i < count; i++ {
		nextPos := Vec2{ pos.x, pos.y - 1 }
		if nextPos.y < col.offset {
			nextPos.y = col.offset + col.size - 1
		}

		switch m.tiles.Get(nextPos.x, nextPos.y) {
			case Wall: return pos
			case Open: pos = nextPos
			case Empty: panic(nextPos)
		}
	}
	return pos
}

func MoveDown(m *Map, pos Vec2, count int) Vec2 {
	// count %= m.rows[pos.y].size
	col := m.cols[pos.x]
	for i := 0; i < count; i++ {
		nextPos := Vec2{ pos.x, pos.y + 1 }
		if nextPos.y >= col.offset + col.size {
			nextPos.y = col.offset
		}

		switch m.tiles.Get(nextPos.x, nextPos.y) {
			case Wall: return pos
			case Open: pos = nextPos
			case Empty: panic(nextPos)
		}
	}
	return pos
}


func MoveLeft(m *Map, pos Vec2, count int) Vec2 {
	// count %= m.rows[pos.y].size
	row := m.rows[pos.y]
	for i := 0; i < count; i++ {
		nextPos := Vec2{ pos.x - 1, pos.y }
		if nextPos.x < row.offset {
			nextPos.x = row.offset + row.size - 1
		}

		switch m.tiles.Get(nextPos.x, nextPos.y) {
			case Wall: return pos
			case Open: pos = nextPos
			case Empty: panic(nextPos)
		}
	}
	return pos
}

func MoveRight(m *Map, pos Vec2, count int) Vec2 {
	// count %= m.rows[pos.y].size
	row := m.rows[pos.y]
	for i := 0; i < count; i++ {
		nextPos := Vec2{ pos.x + 1, pos.y }
		if nextPos.x >= row.offset + row.size {
			nextPos.x = row.offset
		}

		switch m.tiles.Get(nextPos.x, nextPos.y) {
			case Wall: return pos
			case Open: pos = nextPos
			case Empty: panic(nextPos)
		}
	}
	return pos
}

func NewMap(width, height int) Map {
	return Map{
		tiles: NewGrid[Tile](width, height, Empty),
		rows: make([]Span, height),
		cols: make([]Span, width),
	}
}

func getMoveDistance(moves string) (int, string) {
	distance := 0

	for len(moves) > 0 && moves[0] >= '0' && moves[0] <= '9' {
		distance = distance * 10 + int(moves[0] - '0')
		moves = moves[1:]
	}

	return distance, moves
}

func getTurn(moves string) (int, string) {
	turn := moves[0]
	moves = moves[1:]

	//fmt.Printf("Turn %c\n", turn)

	switch turn {
		case 'L': return -1, moves
		case 'R': return +1, moves
		default: panic(turn)
	}
}

type Grid[T any] struct {
    w, h int
	empty T
    cell []T
}

func NewGrid[T any](w, h int, empty T) Grid[T] {
    grid := Grid[T]{
        w: w,
        h: h,
		empty: empty,
        cell: make([]T, w * h),
    }

	for i := 0; i < w * h; i++ {
		grid.cell[i] = empty
	}

    return grid
}

func (this Grid[T]) Get(x, y int) T {
    if this.Contains(x, y) {
        return this.cell[this.offset(x, y)]
    }
    return this.empty
}

func (this *Grid[T]) Set(x, y int, value T) {
    this.cell[this.offset(x, y)] = value
}

func (this Grid[T]) Width() int {
    return this.w
}

func (this Grid[T]) Height() int {
    return this.h
}

func (this Grid[T]) offset(x, y int) int {
    return this.w * y + x
}

func (this Grid[T]) Contains(x, y int) bool {
    return x >= 0 && y >= 0 && x < this.w && y < this.h
}
