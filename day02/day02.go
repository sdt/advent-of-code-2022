package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Move int

const (
	Rock Move = iota
	Paper
	Scissors
)

type Outcome int

const (
	Lose Outcome = iota
	Draw
	Win
)

type Game1 struct {
	opponent Move
	you      Move
}

type Game2 struct {
	opponent Move
	outcome  Outcome
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	total := 0

	for _, line := range lines {
		game := parseGame1(line)
		total += game.score()
	}
	return total
}

func part2(lines []string) int {
	total := 0

	for _, line := range lines {
		game2 := parseGame2(line)
		game1 := game2.solve()
		total += game1.score()
	}
	return total
}

func RPS(rps Move) string {
	switch rps {
	case Rock:
		return "R1"
	case Paper:
		return "P2"
	case Scissors:
		return "S3"
	}
	return "-"
}

func (this Game1) score() int {
	base := this.you.value()
	if this.you == this.opponent {
		return 3 + base
	}
	if this.you.beats() == this.opponent {
		return 6 + base
	}
	return base
}

func (this Game2) solve() Game1 {
	game := Game1{}
	game.opponent = this.opponent

	if this.outcome == Lose {
		game.you = this.opponent.beats()
	} else if this.outcome == Win {
		game.you = this.opponent.losesTo()
	} else {
		game.you = this.opponent
	}
	return game
}

func (this Move) beats() Move {
	return (this - 1 + 3) % 3
}

func (this Move) losesTo() Move {
	return (this + 1) % 3
}

func (this Move) value() int {
	return int(this) + 1
}

func parseGame1(line string) Game1 {
	game := Game1{}
	game.opponent = Move(line[0] - 'A')
	game.you = Move(line[2] - 'X')
	return game
}

func parseGame2(line string) Game2 {
	game := Game2{}
	game.opponent = Move(line[0] - 'A')
	game.outcome = Outcome(line[2] - 'X')
	return game
}
