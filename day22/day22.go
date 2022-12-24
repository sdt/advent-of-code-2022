package main

import (
	"advent-of-code/aoc"
	"fmt"
)

/*
	Top row

	U0 R0 D2 L0 ...
    |  |  |  |
	|  |  |  F1 U1 B1 D1
	|  |  |
	|  |  F2 L1 B0 R3
	|  |
	|  F3 D3 B3 U3
	|
	F0 R1 B2 L3

	   rt dn
	U0 R0 F0
	U1 B1 R1
	U2 L2 B2
	U3 F3 L3

	B0 R3 U0
	B1 D1 R0
	B2 L3 D2
	B3 U3 L0

	D0 R2 B0
	D1 F1 R3
	D2 L0 F2
	D3 B3 L1

	F0 R1 D0
	F1 U1 R2
	F2 L1 U2
	F3 D3 L2

 	L0 U0 F1
	L1 B0 U1
	L2 D0 B1
	L3 F0 D1

	R0 D2 F3
	R1 B2 D3
	R2 U2 B3
	R3 F2 U3
*/

/* Canonical unfold. All sides at rot=0
       +--+
       |F3|
       +--+
       |D2|
       +--+
       |B1|
    +--+--+--+
    |L4|U0|R5|
    +--+--+--+
*/

type CubeSide uint8
const (
	U = iota
	B
	D
	F
	L
	R
)

type Rotation uint8

type CubeRotation struct {
	side CubeSide
	rotation Rotation
}

var rightFace = map[CubeRotation]CubeRotation{
	{U,0}: {R,0}, {U,1}: {B,1}, {U,2}: {L,2}, {U,3}: {F,3},
	{B,0}: {R,3}, {B,1}: {D,1}, {B,2}: {L,3}, {B,3}: {U,3},
	{D,0}: {R,2}, {D,1}: {F,1}, {D,2}: {L,0}, {D,3}: {B,3},
	{F,0}: {R,1}, {F,1}: {U,1}, {F,2}: {L,1}, {F,3}: {D,3},
 	{L,0}: {U,0}, {L,1}: {B,0}, {L,2}: {D,0}, {L,3}: {F,0},
	{R,0}: {D,2}, {R,1}: {B,2}, {R,2}: {U,2}, {R,3}: {F,2},
}

type Vec2 struct {
	x, y int
}

type Tile struct {
	row, col uint8
	isOpen bool
}

type Face struct {
	tile Grid[Tile]
}

type FaceSet struct {
	face map[Vec2]Face
}

// tile size = (X + Y) / 7

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part2(lines))
}

func part2(lines []string) int {
	faceSet := parseFaceSet(lines)

	fmt.Println(faceSet)

	return len(faceSet.face)
}

func parseFaceSet(lines []string) FaceSet {
	h := len(lines)
	w := 0
	for _, line := range lines {
		if len(line) > w {
			w = len(line)
		}
	}

	tileSize := (w + h) / 7

	xTiles := w / tileSize
	yTiles := h / tileSize

	faceSet := FaceSet{ make(map[Vec2]Face) }

	for yTile := 0; yTile < yTiles; yTile++ {
		for xTile := 0; xTile < xTiles; xTile++ {
			yOffset := yTile * tileSize
			xOffset := xTile * tileSize

			if len(lines[yOffset]) <= xOffset {
				continue
			}
			if lines[yOffset][xOffset] == ' ' {
				continue
			}

			face := Face{ NewGrid[Tile](tileSize, tileSize) }
			for y := 0; y < tileSize; y++ {
				for x := 0; x < tileSize; x++ {
					row := uint8(yOffset + y + 1)
					col := uint8(xOffset + x + 1)
					isOpen := lines[yOffset + y][xOffset + x] == '.'
					tile := Tile{ row, col, isOpen }
					face.tile.Set(x, y, tile)
				}
			}

			faceSet.face[ Vec2{ xTile, yTile } ] = face
		}
	}

	return faceSet
}

type Grid[T any] struct {
	w, h int
	cell []T
}

func NewGrid[T any](w, h int) Grid[T] {
	grid := Grid[T]{
		w: w,
		h: h,
		cell: make([]T, w * h),
	}

	return grid
}

func (this Grid[T]) Get(x, y int) T {
	return this.cell[this.offset(x, y)]
}

func (this *Grid[T]) Set(x, y int, value T) {
	this.cell[this.offset(x, y)] = value
}

func (this Grid[T]) offset(x, y int) int {
	return this.w * y + x
}

func (this Grid[T]) Contains(x, y int) bool {
	return x >= 0 && y >= 0 && x < this.w && y < this.h
}

func (this Grid[T]) Rotate(count int) Grid[T] {
	n := count % 4
	if n < 0 {
		n += 4
	}

	switch n {

	case 0: return this.Rotate0()
	case 1: return this.Rotate90()
	case 2: return this.Rotate180()
	case 3: return this.Rotate270()

	}

	panic(n)
}

func (this Grid[T]) Rotate0() Grid[T] {
	return this
}

func (this Grid[T]) Rotate90() Grid[T] {
	that := NewGrid[T](this.h, this.w)

	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			that.Set(that.w - y - 1, x, this.Get(x, y))
		}
	}

	return that
}

func (this Grid[T]) Rotate180() Grid[T] {
	that := NewGrid[T](this.w, this.h)

	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			that.Set(this.w - x - 1, this.h - y - 1, this.Get(x, y))
		}
	}

	return that
}

func (this Grid[T]) Rotate270() Grid[T] {
	that := NewGrid[T](this.h, this.w)

	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			that.Set(y, that.h - x - 1, this.Get(x, y))
		}
	}

	return that
}

func (this Grid[T]) Print(format string) {
	for y, offset := 0, 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			fmt.Printf(format, this.cell[offset])
			offset++
		}
		fmt.Println()
	}
	fmt.Println()
}
