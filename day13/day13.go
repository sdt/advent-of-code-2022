package main

import (
	"advent-of-code/aoc"
	"fmt"
    "sort"
)

type ValueType int
const (
    IntValue  ValueType = iota
    ListValue
)

type Comparison int
const (
    RightOrder Comparison = iota
    WrongOrder
    Indeterminate
)

type PacketValue struct {
    valueType   ValueType
    intValue    int
    listValue   []*PacketValue
}

type Entry struct {
    line    string
    packet  *PacketValue
}

type Entries []Entry

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
            break;
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

    var entries Entries
    for _, line := range lines {
        if len(line) == 0 {
            continue
        }
        entries = append(entries, Entry{ line, parseLine(line) })
    }
    sort.Sort(entries)

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
            valueType: ListValue,
            listValue: listValue,
        }
    }
    rest, intValue := parseIntValue(line)
    return rest, &PacketValue{
        valueType: IntValue,
        intValue:  intValue,
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
        value = value * 10 + int(line[0] - '0')
        line = line[1:]
    }
    //fmt.Printf("['%s', %d]\n", line, value)
    return line, value
}

func (this PacketValue) String() string {
    if this.valueType == IntValue {
        return fmt.Sprintf("%d", this.intValue)
    }

    s := "["
    for i, v := range this.listValue {
        if i != 0 {
            s += ","
        }
        s += v.String()
    }
    s += "]"
    return s
}

func (this PacketValue) MakeListValue() PacketValue {
    if this.valueType == ListValue {
        return this
    }

    return PacketValue{
        valueType: ListValue,
        listValue: []*PacketValue{ &this },
    }
}

func compare(left PacketValue, right PacketValue) Comparison {
    if left.valueType == IntValue && right.valueType == IntValue {
        switch {
            case left.intValue < right.intValue:
                return RightOrder

            case left.intValue > right.intValue:
                return WrongOrder

            default:
                return Indeterminate
        }
    }

    if left.valueType == ListValue && right.valueType == ListValue {
        leftVal := left.listValue
        rightVal := right.listValue

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

func (this Entries) Len() int {
    return len(this)
}

func (this Entries) Less(i, j int) bool {
    return compare(*this[i].packet, *this[j].packet) == RightOrder
}

func (this Entries) Swap(i, j int) {
    this[i], this[j] = this[j], this[i]
}
