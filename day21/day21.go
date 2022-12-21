package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

type Operator string

type Operation struct {
	lhs, rhs string
	operator Operator
}

type Monkey struct {
	job any	// Operation or int
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	monkeyTable := make(map[string]Monkey, len(lines))

	for _, line := range lines {
		name, monkey := parseMonkey(line)
		monkeyTable[name] = monkey
	}

	var evaluate func(string) int
	evaluate = func(name string) int {
		value := monkeyTable[name]
		if intValue, isIntValue := value.job.(int); isIntValue {
			return intValue
		}

		operation := value.job.(Operation)

		result := doOperator(
			operation.operator,
			evaluate(operation.lhs),
			evaluate(operation.rhs),
		)

		monkeyTable[name] = Monkey{ job: result }

		return result
	}
	return evaluate("root")
}

func doOperator(operator Operator, lhs, rhs int) int {
	switch (operator) {

	case "+": return lhs + rhs
	case "-": return lhs - rhs
	case "*": return lhs * rhs
	case "/": return lhs / rhs
	default: panic(operator)

	}
}

func parseMonkey(line string) (string, Monkey) {
	words := strings.Split(line, " ")
	name := words[0]
	name = name[:len(name)-1] // drop the colon

	if len(words) == 2 {
		return name, Monkey{ job: aoc.ParseInt(words[1]) }
	}

	return name, Monkey{
		job: Operation{
			lhs: words[1],
			rhs: words[3],
			operator: Operator(words[2]),
		},
	}
}
