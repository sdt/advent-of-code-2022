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

var sideName = [...]string{ "U", "B", "D", "F", "L", "R" }

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

var downFace = map[CubeRotation]CubeRotation{
	{U,0}: {F,0}, {U,1}: {R,1}, {U,2}: {B,2}, {U,3}: {L,3},
	{B,0}: {U,0}, {B,1}: {R,0}, {B,2}: {D,2}, {B,3}: {L,0},
	{D,0}: {B,0}, {D,1}: {R,3}, {D,2}: {F,2}, {D,3}: {L,1},
	{F,0}: {D,0}, {F,1}: {R,2}, {F,2}: {U,2}, {F,3}: {L,2},
	{L,0}: {F,1}, {L,1}: {U,1}, {L,2}: {B,1}, {L,3}: {D,1},
	{R,0}: {F,3}, {R,1}: {D,3}, {R,2}: {B,3}, {R,3}: {U,3},
}

// leftFace is the reverse mapping of rightface
var leftFace = ReverseMap(rightFace)

type Vec2 struct {
	x, y int
}

func (this Vec2) Add(that Vec2) Vec2 {
	return Vec2{ this.x + that.x, this.y + that.y }
}

type Tile struct {
	row, col uint8
	isOpen bool
}

type Face struct {
	tile Grid[Tile]
}

type FaceSet struct {
	face map[Vec2]*Face
}

// The cube is represented as an array of six faces.
// Each face has four rotations precomputed.
// The 0 rotation matches the canonical representation. The rest are the CW
// rotations from the previous.
type Cube struct {
	face [6][4]*Face
}

type CubeOrientation struct {
	cubeRot [6]CubeRotation
}

// tile size = (X + Y) / 7

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part2(lines))
}

func part2(lines []string) int {
	// We start facing up, so always insert an R
	cube := buildCube(parseFaceSet(lines[:len(lines)-2]))
	orientation := NewCubeOrientation()
	moves := "R" + lines[len(lines)-1]
	pos := Vec2{ 0, 0 }
	size := cube.face[0][0].tile.w

	DrawCube(&cube, &orientation)

	GetTile := func(p Vec2, face int) Tile {
		cr := orientation.cubeRot[face]
		f := cube.face[cr.side][cr.rotation]
		return f.tile.Get(p.x, p.y)
	}

	for len(moves) > 0 {
		turn, nextMoves := getTurn(moves)
		moves = nextMoves

		switch turn {
		case 'R':
			orientation = TurnRight(orientation)
			// 1,2 -> 2,S-1-1
			pos = Vec2{ pos.y, size - pos.x - 1 }

		case 'L':
			orientation = TurnLeft(orientation)
			// 1,2 -> S-1-1, 2
			pos = Vec2{ size - pos.y - 1, pos.x }

		default: panic(turn)
		}
		fmt.Printf("Turn %c: %d,%d\n", turn, pos.x, pos.y)
		DrawCube(&cube, &orientation)


		distance, nextMoves := getMoveDistance(moves)
		fmt.Println("Move:", distance)
		moves = nextMoves
		for i := 0; i < distance; i++ {
			fmt.Printf("row,col=%d,%d\n", GetTile(pos, 0).row, GetTile(pos, 0).col)
			if pos.y > 0 {
				nextPos := pos.Add(Vec2{ 0, -1 })
				if !GetTile(nextPos, 0).isOpen {
					fmt.Printf("Step %d of %d: blocked at %d,%d\n", i+1, distance, pos.x, pos.y)
					break
				}
				pos = nextPos
				fmt.Printf("Step %d of %d: %d,%d\n", i+1, distance, pos.x, pos.y)
				DrawCube(&cube, &orientation)
			} else {
				nextPos := Vec2{ pos.x, size - 1 }
				if !GetTile(nextPos, 1).isOpen {
					fmt.Printf("Step %d of %d: blocked around corner at %d,%d\n", i+1, distance, pos.x, pos.y)
					break
				}
				orientation = RollForward(orientation)
				pos = nextPos
				fmt.Printf("Jump %d of %d: %d,%d\n", i+1, distance, pos.x, pos.y)
				DrawCube(&cube, &orientation)
			}
		}
	}

	final := GetTile(pos, 0)
	fmt.Printf("Final position: face=%d row=%d,col=%d dir=%d\n", orientation.cubeRot[0].side, final.row, final.col, (orientation.cubeRot[0].rotation + 0) % 4)

	return int(final.row) * 1000 +
		   int(final.col) * 4 +
		   int(orientation.cubeRot[0].rotation)
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

	faceSet := FaceSet{ make(map[Vec2]*Face) }

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

			face := NewFace(tileSize)
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

func buildCube(faceSet FaceSet) Cube {

	var cube Cube

	// First, find the up face. It will be the lowest X,0.
	var up Vec2
	for x := 0; ; x++ {
		up = Vec2{x, 0}
		if face, found := faceSet.face[up]; found {
			delete(faceSet.face, up)
			fmt.Printf("Face %s is %d,%d at rotation %d\n", sideName[U], up.x, up.y, 0)
			setFace(&cube, face, CubeRotation{U,0})
			break
		}
	}

	var flood func(Vec2, CubeRotation)
	flood = func(pos Vec2, rot CubeRotation) {

		var nextPos Vec2
		var nextRot CubeRotation

		// Try the square to the right
		nextPos = pos.Add(Vec2{1, 0})
		if face, found := faceSet.face[nextPos]; found {
			nextRot = rightFace[rot]
			delete(faceSet.face, nextPos)
			setFace(&cube, face, nextRot)
			fmt.Printf("Face %s is %d,%d at rotation %d\n", sideName[nextRot.side], nextPos.x, nextPos.y, nextRot.rotation)
			flood(nextPos, nextRot)
		}

		// Try the square to the left
		nextPos = pos.Add(Vec2{-1, 0})
		if face, found := faceSet.face[nextPos]; found {
			nextRot = leftFace[rot]
			delete(faceSet.face, nextPos)
			setFace(&cube, face, nextRot)
			fmt.Printf("Face %s is %d,%d at rotation %d\n", sideName[nextRot.side], nextPos.x, nextPos.y, nextRot.rotation)
			flood(nextPos, nextRot)
		}

		// Try the square down
		nextPos = pos.Add(Vec2{0, 1})
		if face, found := faceSet.face[nextPos]; found {
			nextRot = downFace[rot]
			delete(faceSet.face, nextPos)
			setFace(&cube, face, nextRot)
			fmt.Printf("Face %s is %d,%d at rotation %d\n", sideName[nextRot.side], nextPos.x, nextPos.y, nextRot.rotation)
			flood(nextPos, nextRot)
		}


	}
	flood(up, CubeRotation{U,0})

	return cube
}

func setFace(cube *Cube, face* Face, cubeRot CubeRotation) {
	for i := 0; i < 4; i++ {
		cube.face[cubeRot.side][(int(cubeRot.rotation) + i) % 4] = face
		face = rotateFace(face)
	}
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

func NewFace(size int) *Face {
	return &Face{ NewGrid[Tile](size, size) }
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

func rotateFace(face* Face) *Face {
	return &Face{ face.tile.RotateCW() }
}

func (this Grid[T]) RotateCW() Grid[T] {
	that := NewGrid[T](this.h, this.w)

	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			that.Set(that.w - y - 1, x, this.Get(x, y))
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

func ReverseMap[T comparable](from map[T]T) map[T]T {
	to := make(map[T]T)

	for k, v := range from {
		to[v] = k
	}
	return to
}

func NewCubeOrientation() CubeOrientation {
	return CubeOrientation{ [...]CubeRotation{
		{U, 0}, {B, 0}, {D, 0},
		{F, 0}, {L, 0}, {R, 0},
	} }
}

func RollForward(orientation CubeOrientation) CubeOrientation {
	cr := orientation.cubeRot;
	return CubeOrientation{ [...]CubeRotation{
		cr[1], // 1 -> 0
		cr[2], // 2 -> 1
		cr[3], // 3 -> 2
		cr[0], // 0 -> 3
		{ cr[4].side, (cr[4].rotation + 1) % 4 },
		{ cr[5].side, (cr[5].rotation + 3) % 4 },
	}}
}

func TurnRight(orientation CubeOrientation) CubeOrientation {
	cr := orientation.cubeRot
	side := func(index, rot int) CubeRotation {
		return CubeRotation{
			cr[index].side,
			Rotation((int(cr[index].rotation) + rot + 4) % 4),
		}
	}
	return CubeOrientation{ [...]CubeRotation{
		side(0, -1),
		side(5, -1),
		side(2,  1),
		side(4, -1),
		side(1, -1),
		side(3, -1),
	}}
}

func TurnLeft(orientation CubeOrientation) CubeOrientation {
	cr := orientation.cubeRot
	side := func(index, rot int) CubeRotation {
		return CubeRotation{
			cr[index].side,
			Rotation((int(cr[index].rotation) + rot + 4) % 4),
		}
	}
	return CubeOrientation{ [...]CubeRotation{
		side(0,  1),
		side(4,  1),
		side(2, -1),
		side(5,  1),
		side(3,  1),
		side(1,  1),
	}}
}

func DrawCube(cube *Cube, orientation* CubeOrientation) {
	return
	size := cube.face[0][0].tile.w

	drawRow := func(face* Face, y int) {
		for x := 0; x < size; x++ {
			tile := face.tile.Get(x, y)
			if tile.isOpen {
				fmt.Printf("%c", '.')
			} else {
				fmt.Printf("%c", '#')
			}
		}
	}

	drawGap := func() {
		for x := 0; x < size+1; x++ {
			fmt.Printf("%c", ' ')
		}
	}

	for i := 3; i >= 1; i-- {
		cubeRot := orientation.cubeRot[i]
		face := cube.face[cubeRot.side][cubeRot.rotation]

		for row := 0; row < size; row++ {
			drawGap()
			drawRow(face, row)
			fmt.Println()
		}
		fmt.Println()
	}

	faces := [...]int{4, 0, 5}
	for row := 0; row < size; row++ {
		for _, i := range faces {
			cubeRot := orientation.cubeRot[i]
			face := cube.face[cubeRot.side][cubeRot.rotation]

			drawRow(face, row)
			fmt.Print(" ")
		}
		fmt.Println()
	}
	fmt.Println()
}


func getMoveDistance(moves string) (int, string) {
	distance := 0
	for len(moves) > 0 && moves[0] >= '0' && moves[0] <= '9' {
		distance = distance * 10 + int(moves[0] - '0')
		moves = moves[1:]
	}
	return distance, moves
}

func getTurn(moves string) (byte, string) {
	turn := moves[0]
	moves = moves[1:]

	return turn, moves
}

// 142383 too high
// 142380 correct, but hmmmm
// 142381
// 142382
