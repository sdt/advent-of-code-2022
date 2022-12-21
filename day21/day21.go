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

type Value any

type Expression struct {
	lhs      Value
	operator Operator
	rhs      Value
}

type Monkey any // Operation or int

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
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
		if intValue, isIntValue := value.(int); isIntValue {
			return intValue
		}

		operation := value.(Operation)

		result := doOperator(
			operation.operator,
			evaluate(operation.lhs),
			evaluate(operation.rhs),
		)

		monkeyTable[name] = Monkey(result)

		return result
	}
	return evaluate("root")
}

func part2(lines []string) int {
	monkeyTable := make(map[string]Monkey, len(lines))

	var root Operation

	for _, line := range lines {
		name, monkey := parseMonkey(line)
		if name == "root" {
			root = monkey.(Operation)
		} else if name != "humn" {
			monkeyTable[name] = monkey
		}
	}

	var evaluate func(string) Value
	evaluate = func(name string) Value {
		monkey := monkeyTable[name]
		if monkey == nil {
			return Value(name)
		}
		if intValue, isIntValue := monkey.(int); isIntValue {
			return Value(intValue)
		}

		operation := monkey.(Operation)

		lhs := evaluate(operation.lhs)
		rhs := evaluate(operation.rhs)

		if intLhs, isIntLhs := lhs.(int); isIntLhs {
			if intRhs, isIntRhs := rhs.(int); isIntRhs {
				return Value(doOperator(operation.operator, intLhs, intRhs))
			}
		}

		return Value(Expression{ lhs, operation.operator, rhs })
	}

	var solve func(Value, int) int
	solve = func(lhs Value, rhs int) int {
		//fmt.Println("Solve", lhs, rhs)
		if _, isString := lhs.(string); isString {
			return rhs
		}

		exp := lhs.(Expression)

		if innerLhs, isInt := exp.lhs.(int); isInt {
			subValue := exp.rhs

			// A op x = B
			switch exp.operator {

			// A + x = B =>  x = B - A
			case "+": return solve(subValue, rhs - innerLhs)

			// A - x = B => x = A - B
			case "-": return solve(subValue, innerLhs - rhs)

			// A * x = B =>  x = B / A
			case "*": return solve(subValue, rhs / innerLhs)

			// A / x = B => x = A / B
			case "/": return solve(subValue, innerLhs / rhs)

			}
		} else {
			innerRhs := exp.rhs.(int) // A
			subValue := exp.lhs

			// x op A = B
			switch exp.operator {

			// x + A = B =>  x = B - A
			case "+": return solve(subValue, rhs - innerRhs)

			// x - A = B => x = A + B
			case "-": return solve(subValue, innerRhs + rhs)

			// a * A = B =>  x = B / A
			case "*": return solve(subValue, rhs / innerRhs)

			// x / A = B => x = A * B
			case "/": return solve(subValue, innerRhs * rhs)

			}
		}

		panic(lhs)
	}

	lhs := evaluate(root.lhs)
	rhs := evaluate(root.rhs)

	if lhsValue, lhsIsInt := lhs.(int); lhsIsInt {
		return solve(rhs.(Expression), lhsValue)
	}

	return solve(lhs.(Expression), rhs.(int))
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
		return name, Monkey(aoc.ParseInt(words[1]))
	}

	return name, Monkey(Operation{
		lhs: words[1],
		rhs: words[3],
		operator: Operator(words[2]),
	})
}
