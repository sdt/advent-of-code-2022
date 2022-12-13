package main

import (
	"advent-of-code/aoc"
	"fmt"
	"sort"
	"strings"
)

type Op func(int) int

type Monkey struct {
	items                   []int
	op                      Op
	mod                     int
	div                     int
	test                    int
	trueMonkey, falseMonkey int
	inspections             int
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	monkeys := make([]Monkey, 0)

	for i := 0; i < len(lines); i += 7 {
		monkey := parseMonkey(lines[i:])
		monkey.div = 3
		monkeys = append(monkeys, monkey)
	}

	for i := 0; i < 20; i++ {
		doRound(monkeys)
	}

	inspections := make([]int, len(monkeys))
	for i := range monkeys {
		inspections[i] = monkeys[i].inspections
	}
	sort.Sort(sort.Reverse(sort.IntSlice(inspections)))

	return inspections[0] * inspections[1]
}

func part2(lines []string) int {
	monkeys := make([]Monkey, 0)

	for i := 0; i < len(lines); i += 7 {
		monkey := parseMonkey(lines[i:])
		monkeys = append(monkeys, monkey)
	}

	totalMod := 1
	for i := range monkeys {
		totalMod *= monkeys[i].test
	}
	for i := range monkeys {
		monkeys[i].mod = totalMod
	}

	for i := 0; i < 10000; i++ {
		doRound(monkeys)
	}

	inspections := make([]int, len(monkeys))
	for i := range monkeys {
		inspections[i] = monkeys[i].inspections
	}
	sort.Sort(sort.Reverse(sort.IntSlice(inspections)))

	return inspections[0] * inspections[1]
}

func doRound(monkeys []Monkey) {
	for i := range monkeys {
		for _, item := range monkeys[i].items {
			newItem, next := monkeys[i].doItem(item)
			monkeys[next].items = append(monkeys[next].items, newItem)
		}
		monkeys[i].inspections += len(monkeys[i].items)
		monkeys[i].items = monkeys[i].items[:0]
	}
}

func (this Monkey) doItem(item int) (int, int) {
	item = this.op(item%this.mod) / this.div

	if int(item)%this.test == 0 {
		return item, this.trueMonkey
	}
	return item, this.falseMonkey
}

func parseMonkey(lines []string) Monkey {
	return Monkey{
		items:       parseItems(lines[1]),
		op:          parseOp(lines[2]),
		mod:         10000000000,
		div:         1,
		test:        parseInt(lines[3], "  Test: divisible by "),
		trueMonkey:  parseInt(lines[4], "    If true: throw to monkey "),
		falseMonkey: parseInt(lines[5], "    If false: throw to monkey "),
	}
}

func parseItems(line string) []int {
	prefix := "  Starting items: "
	line = line[len(prefix):]

	tokens := strings.Split(line, ", ")
	return aoc.ParseInts(tokens)
}

func parseOp(line string) Op {
	prefix := "  Operation: new = old "
	line = line[len(prefix):]

	if line == "* old" {
		return func(x int) int { return x * x }
	}
	tokens := strings.Split(line, " ")
	arg := aoc.ParseInt(tokens[1])

	if tokens[0] == "+" {
		return func(x int) int { return x + arg }
	}
	return func(x int) int { return x * arg }
}

func parseInt(line string, prefix string) int {
	return aoc.ParseInt(line[len(prefix):])
}
