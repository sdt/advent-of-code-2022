package main


// 3178 too low

import (
	"advent-of-code/aoc"
	"bufio"
	"fmt"
	"os"
	"time"
)

const debug = !true

type Vec2 struct {
	x, y int
}

type Piece []Vec2

type Chamber struct {
	cell map[Vec2] bool
	height int
}

func main() {
	filename := aoc.GetFilename()
	input := aoc.Slurp(filename)

	for {
		last := len(input)-1
		if input[last] == '<' || input[last] == '>' {
			break
		}
		input = input[:last]
	}

	fmt.Println(part1(input))
}

func part1(input string) int {
	pieces := makePieces()
	chamber := makeChamber()

	move := 0

	for i := 0; i < 2022; i++ {
		piece := pieces[i % len(pieces)]
		move = chamber.dropPiece(&piece, input, move)
		if debug {
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}

	return chamber.height
}

func makePieces() []Piece {
	return []Piece{
		Piece{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		Piece{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},
		Piece{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
		Piece{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		Piece{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	}
}

func makeChamber() *Chamber {
	return &Chamber{ make(map[Vec2]bool), 0 }
}

func (this *Chamber) dropPiece(piece *Piece, moves string, move int) int {
	down := Vec2{0, -1}
	offset := Vec2{2, this.height + 3}
	this.printWith(piece, offset)
	for {
		wind := delta(moves[move])
		move = (move + 1) % len(moves)
		if this.canPlace(piece, offset.add(wind)) {
			offset = offset.add(wind)
			this.printWith(piece, offset)
		}

		if !this.canPlace(piece, offset.add(down)) {
			break
		}
		offset = offset.add(down)
		this.printWith(piece, offset)
	}
	this.place(piece, offset)
	return move
}

func (this *Chamber) isClear(p Vec2) bool {
	if (p.y < 0) || (p.x < 0) || (p.x > 6) {
		return false
	}

	_, found := this.cell[p]
	return !found
}

func (this *Chamber) canPlace(piece *Piece, offset Vec2) bool {
	for _, p := range *piece {
		if !this.isClear(p.add(offset)) {
			return false
		}
	}
	return true
}

func (this *Chamber) place(piece *Piece, offset Vec2) {
	for _, p := range *piece {
		p = p.add(offset)
		if p.y+1 > this.height {
			this.height = p.y+1
		}
		this.cell[p] = true
	}
}

func (this *Chamber) printWith(piece *Piece, offset Vec2) {
	if !debug {
		return
	}
	fmt.Print("\x1b[H\x1b[J")
	fmt.Println(this.height)
	extra := make(map[Vec2]bool)
	for _, p := range *piece {
		p = p.add(offset)
		extra[p] = true
	}

	for y := this.height + 5; y >= 0; y-- {
		fmt.Print("#")
		for x := 0; x <= 6; x++ {
			v := Vec2{x, y}
			if _, onPiece := extra[v]; onPiece {
				fmt.Print("O")
			} else if this.isClear(Vec2{x, y}) {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println("#")

		if this.height - y > 20 {
			fmt.Println("#~~~~~~~#\n")
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
	fmt.Println("#########\n")
	time.Sleep(100 * time.Millisecond)
}

func (this *Chamber) print() {
	fmt.Print("\x1b[H\x1b[J")
	fmt.Println(this.height)
	for y := this.height; y >= 0; y-- {
		fmt.Print("#")
		for x := 0; x <= 6; x++ {
			if this.isClear(Vec2{x, y}) {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println("#")

		if this.height - y > 20 {
			fmt.Println("#~~~~~~~#\n")
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
	fmt.Println("#########\n")
	time.Sleep(500 * time.Millisecond)
}

func (this Vec2) add(that Vec2) Vec2 {
	return Vec2{this.x + that.x, this.y + that.y}
}

func delta(move byte) Vec2 {
	switch move {

	case '<': return Vec2{-1, 0}
	case '>': return Vec2{+1, 0}

	default: panic(move)
	}
}
