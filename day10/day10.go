package main

import (
    "advent-of-code/aoc"
    "fmt"
    "strings"
)

type VM struct {
    x int
    cycle int
    signalCycles []int
}

func main() {
    filename := aoc.GetFilename()
    lines := aoc.GetInputLines(filename)

    fmt.Println(part1(lines))
}

func part1(lines []string) int {
    vm := NewVM()
    signalCycles := []int{ 20, 60, 100, 140, 180, 220 };

    signalStrength := 0

    for _, line := range lines {
        incr := 0
        cycles := 1

        tokens := strings.Split(line, " ")
        if tokens[0] == "addx" {
            incr = aoc.ParseInt(tokens[1])
            cycles = 2
        }

        for i := 0; i < cycles; i++ {
            vm.cycle++
            if (len(signalCycles) > 0) && (signalCycles[0] == vm.cycle) {
                signalCycles = signalCycles[1:]
                signalStrength += vm.x * vm.cycle
            }
        }
        vm.x += incr
    }
    return signalStrength
}

func NewVM(signalCycles ...int) VM {
    return VM{ x: 1, cycle: 0, signalCycles: signalCycles }
}
