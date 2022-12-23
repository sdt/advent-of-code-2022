package main

import (
	"advent-of-code/aoc"
	"fmt"
)

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
