package main

import (
	"advent-of-code/aoc"
	"fmt"
	"sort"
)

type ValueType int

type Comparison int

const (
	RightOrder Comparison = iota
	WrongOrder
	Indeterminate
)

type PacketValue struct {
	value any // either an int, or an array of PacketValues
}

type Entry struct {
	line   string
	packet *PacketValue
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	pair := 1
	result := 0
	for {
		left := parseLine(lines[0])
		right := parseLine(lines[1])
		comp := compare(*left, *right)
		//fmt.Println(pair, left, right, comp)

		if comp == RightOrder {
			result += pair
		}

		if len(lines) == 2 {
			break
		}
		lines = lines[3:]
		pair++
	}
	return result
}

func part2(lines []string) int {
	first := "[[2]]"
	second := "[[6]]"
	lines = append(lines, first, second)

	var entries []Entry
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		entries = append(entries, Entry{line, parseLine(line)})
	}
	sort.Slice(entries, func(i, j int) bool {
		return compare(*entries[i].packet, *entries[j].packet) == RightOrder
	})

	ret := 1
	for i, entry := range entries {
		if entry.line == first || entry.line == second {
			ret = ret * (i + 1)
		}
	}

	return ret
}

func parseLine(line string) *PacketValue {
	rest, value := parse(line)
	if len(rest) != 0 {
		panic("Parse error")
	}
	return value
}

func parse(line string) (string, *PacketValue) {
	//fmt.Printf("parse('%s')\n", line)
	if line[0] == '[' {
		rest, listValue := parseListValue(line[1:])
		return rest, &PacketValue{
			value: listValue,
		}
	}
	rest, intValue := parseIntValue(line)
	return rest, &PacketValue{
		value: intValue,
	}
}

func parseListValue(line string) (string, []*PacketValue) {
	//fmt.Printf("parseList('%s')\n", line)
	values := make([]*PacketValue, 0)
	for line[0] != ']' {
		next, value := parse(line)
		values = append(values, value)
		if next[0] == ',' {
			next = next[1:]
		}
		line = next
	}
	return line[1:], values
}

func parseIntValue(line string) (string, int) {
	//fmt.Printf("parseInt('%s') -> ", line)
	value := 0
	for line[0] >= '0' && line[0] <= '9' {
		value = value*10 + int(line[0]-'0')
		line = line[1:]
	}
	//fmt.Printf("['%s', %d]\n", line, value)
	return line, value
}

func (this PacketValue) String() string {
	if value, ok := this.value.(int); ok {
		return fmt.Sprintf("%d", value)
	}

	s := "["
	for i, v := range this.value.([]*PacketValue) {
		if i != 0 {
			s += ","
		}
		s += v.String()
	}
	s += "]"
	return s
}

func (this PacketValue) MakeListValue() PacketValue {
	if _, ok := this.value.(int); ok {
		return PacketValue{[]*PacketValue{&this}}
	}

	return this
}

func compare(left PacketValue, right PacketValue) Comparison {
	leftInt, leftIsInt := left.value.(int)
	rightInt, rightIsInt := right.value.(int)

	if leftIsInt && rightIsInt {
		switch {
		case leftInt < rightInt:
			return RightOrder

		case leftInt > rightInt:
			return WrongOrder

		default:
			return Indeterminate
		}
	}

	if !leftIsInt && !rightIsInt {
		leftVal := left.value.([]*PacketValue)
		rightVal := right.value.([]*PacketValue)

		for {
			if len(leftVal) == 0 {
				if len(rightVal) == 0 {
					return Indeterminate
				}
				return RightOrder
			}
			if len(rightVal) == 0 {
				return WrongOrder
			}
			comp := compare(*leftVal[0], *rightVal[0])
			if comp != Indeterminate {
				return comp
			}
			leftVal = leftVal[1:]
			rightVal = rightVal[1:]
		}
	}

	return compare(left.MakeListValue(), right.MakeListValue())
}
