package main

import (
	"advent-of-code/aoc"
	"fmt"
    "strings"
)

type Crate rune
type Stack []Crate

type Move struct {
    from, to, howMany int
}

const debug = true

func main() {
	filename := aoc.GetFilename()
	input := aoc.Slurp(filename)

	fmt.Println(part1(input))
}

func part1(input string) string {
    stacks, moves := parseInput(input)

    printStacks(stacks)

    for _, move := range moves {
        stacks = applyMove(stacks, move)
        printStacks(stacks)
    }

    ret := ""
    for _, stack := range stacks {
        ret = fmt.Sprintf("%s%c", ret, stack[len(stack)-1])
    }
    return ret
}

func applyMove(stacks []Stack, move Move) []Stack {
    if debug {
        fmt.Printf("Move %d from %d to %d\n", move.howMany, move.from, move.to);
    }
    from := stacks[move.from - 1]
    to   := stacks[move.to - 1]

    for i := 0; i < move.howMany; i++ {
        to = append(to, from[len(from) - i - 1])
    }

    stacks[move.to - 1] = to
    stacks[move.from - 1] = from[0:len(from) - move.howMany]

    return stacks
}

func parseInput(input string) ([]Stack, []Move) {
    parts := strings.Split(input, "\n\n")

    return parseStacks(parts[0]), parseMoves(parts[1])
}

func parseStacks(input string) []Stack {
    lines := strings.Split(input, "\n")

    counts := lines[len(lines)-1]
    count := (len(counts) + 2) / 4

    stacks := make([]Stack, count)

    for i := len(lines) - 2; i >= 0; i-- {
        line := lines[i]
        j := 1
        for s, _ := range stacks {
            if line[j] != ' ' {
                stacks[s] = append(stacks[s], Crate(line[j]))
            }
            j += 4
        }
    }

    return stacks
}

func printStacks(stacks []Stack) {
    if !debug {
        return
    }

    for i, stack := range stacks {
        fmt.Printf("[%d]", i+1);
        for _, crate := range stack {
            fmt.Printf(" %c", crate);
        }
        fmt.Println("");
    }
    fmt.Println();
}

func parseMoves(input string) []Move {
    moves := make([]Move, 0)

    lines := strings.Split(input, "\n");
    for _, line := range lines {
        if line != "" {
            moves = append(moves, parseMove(line))
        }
    }

    return moves
}

func parseMove(line string) Move {
    fmt.Println("move:", line)
    words := strings.Split(line, " ");
    return Move{
        howMany: aoc.ParseInt(words[1]),
        from:    aoc.ParseInt(words[3]),
        to:     aoc.ParseInt(words[5]),
    }
}
