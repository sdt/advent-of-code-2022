package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
	"strings"
)

var InputRegex = regexp.MustCompile("^Valve ([A-Z]{2}) has flow rate=(\\d+); tunnels? leads? to valves? ([A-Z]{2}(, [A-Z]{2})*)$")

type Valve struct {
	name string
	bit uint64
	rate int
	leadsTo []*Valve
}

type State struct {
	valve *Valve
	openValves uint64
	time int
	rate int
	pressure int
}

type StateSet map[string]*State

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	aa := parseInput(lines)

	current := make(StateSet)

	AddState(current, NewState(aa))

	endTime := 30

	for time := 0; time < endTime-1; time++ {
		next := make(StateSet)

		for _, state := range current {
			AddState(next, state.OpenValve())

			for _, valve := range state.valve.leadsTo {
				AddState(next, state.MoveTo(valve))
			}
		}

		current = next
	}

	max := 0
	for _, state := range current {
		finalPressure := state.pressure + state.rate
		if finalPressure > max {
			max = finalPressure
		}
	}

	return max
}

func AddState(set StateSet, state* State) {
	if state == nil {
		return
	}
	key := state.Key()

	if existing, found := set[key]; found {
		if existing.pressure >= state.pressure {
			return
		}
	}
	set[key] = state
}

func parseInput(lines []string) *Valve {
	byName := make(map[string]*Valve)
	destinations := make(map[string][]string)

	for i, line := range lines {
		name, rate, leadsTo := parseLine(line)
		byName[name] = &Valve{name, 1 << i, rate, make([]*Valve, len(leadsTo))}
		destinations[name] = leadsTo
	}

	for name, leadsTo := range destinations {
		from := byName[name]
		for i, valveName := range leadsTo {
			from.leadsTo[i] = byName[valveName]
		}
	}

	return byName["AA"]
}

func parseLine(line string) (string, int, []string) {
	matches := InputRegex.FindStringSubmatch(line)

	name := matches[1]
	rate := aoc.ParseInt(matches[2])
	leadsTo := strings.Split(matches[3], ", ")

	return name, rate, leadsTo
}

func NewState(valve *Valve) *State {
	return &State{ valve, 0, 0, 0, 0 }
}

func (this *State) MoveTo(next *Valve) *State {
	openValves := this.openValves
	time := this.time + 1
	pressure := this.pressure + this.rate

	return &State{ next, openValves, time, this.rate, pressure }
}

func (this *State) OpenValve() *State {
	if this.valve.rate == 0 {
		return nil // no point opening a zero-flow valve
	}

	if this.openValves & this.valve.bit == this.valve.bit {
		return nil // already open
	}

	openValves := this.openValves | this.valve.bit
	time := this.time + 1
	pressure := this.pressure + this.rate
	rate := this.rate + this.valve.rate

	return &State{ this.valve, openValves, time, rate, pressure }
}

func (this *State) Key() string {
	return fmt.Sprintf("%s/%x", this.valve.name, this.openValves)
}
