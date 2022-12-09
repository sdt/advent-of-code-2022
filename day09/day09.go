package main

import (
    "advent-of-code/aoc"
    "fmt"
    "strings"
)

type Vec2 struct {
    x, y int
}

func main() {
    filename := aoc.GetFilename()
    lines := aoc.GetInputLines(filename)

    fmt.Println(part1(lines))
    fmt.Println(part2(lines))
}

func part1(lines []string) int {
    return run(lines, 2)
}

func part2(lines []string) int {
    return run(lines, 10)
}

func run(lines []string, knots int) int {
    knot := make([]Vec2, knots)

    tailLog := make(map[Vec2]bool)

    tailLog[knot[knots-1]] = true

    for _, line := range lines {
        dir, dist := parseLine(line)

        for i := 0; i < dist; i++ {
            knot[0] = knot[0].Add(dir)
            for j := 0; j < knots - 1; j++ {
                diff := knot[j].Sub(knot[j+1])
                move := diff.Sign()
                if diff != move {
                    knot[j+1] = knot[j+1].Add(move)
                }
            }
            tailLog[knot[knots-1]] = true
        }
    }
    return len(tailLog)
}

func parseLine(line string) (Vec2, int) {
    words := strings.Split(line, " ")

    dist := aoc.ParseInt(words[1])
    var dir Vec2

    switch words[0] {

    case "R": dir = NewVec2(+1, 0);
    case "L": dir = NewVec2(-1, 0);
    case "U": dir = NewVec2(0, +1);
    case "D": dir = NewVec2(0, -1);

    }

    return dir, dist
}

func NewVec2(x, y int) Vec2 {
    return Vec2{x, y}
}

func sign(x int) int {
    if x < 0 {
        return -1
    }
    if x > 0 {
        return +1
    }
    return 0
}

func (this Vec2) Add(that Vec2) Vec2 {
    return NewVec2(this.x + that.x, this.y + that.y)
}

func (this Vec2) Sub(that Vec2) Vec2 {
    return NewVec2(this.x - that.x, this.y - that.y)
}

func (this Vec2) Sign() Vec2 {
    return NewVec2(sign(this.x), sign(this.y))
}
