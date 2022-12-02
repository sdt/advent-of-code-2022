package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type RockPaperScissors int
const (
    Rock RockPaperScissors = 1
    Paper                  = 2
    Scissors               = 3
)

type Game struct {
    you         RockPaperScissors
    opponent    RockPaperScissors
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
    total := 0

	for _, line := range lines {
        game := parseGame(line)
        //fmt.Println(RPS(game.opponent), RPS(game.you), game.score())
        total += game.score()
	}
	return total
}

func RPS(rps RockPaperScissors) string {
    switch rps {
    case Rock: return "R1"
    case Paper: return "P2"
    case Scissors: return "S3"
    }
    return "-"
}

func (this Game) score() int {
    if this.you == this.opponent {
        return 3 + int(this.you)
    }
    if this.you.beats(this.opponent) {
        return 6 + int(this.you)
    }
    return int(this.you)
}

func (this RockPaperScissors) beats (that RockPaperScissors) bool {
    switch this {
    case Rock: return that == Scissors
    case Paper: return that == Rock
    case Scissors: return that == Paper
    }
    return false
}

func parseGame(line string) Game {
    game := Game{}
    game.opponent = parseRockPaperScissors(line[0])
    game.you      = parseRockPaperScissors(line[2])
    return game
}

func parseRockPaperScissors(c byte) RockPaperScissors {
    switch c {
    case 'A': return Rock
    case 'B': return Paper
    case 'C': return Scissors
    case 'X': return Rock
    case 'Y': return Paper
    case 'Z': return Scissors
    }
    return Rock
}
