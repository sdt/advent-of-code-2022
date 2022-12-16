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

type Set[T comparable] struct {
	element map[T]bool
}

type StateQueue struct {
	queue []*State
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	aa := parseInput(lines)

	endTime := 30

	max := 0
	agenda := []*State{ NewState(aa) }

	best := make(map[string]int)

	for len(agenda) > 0 {
		state := agenda[0]
		//fmt.Printf("%s time=%d pressure=%d\n", state.valve.name, state.time, state.pressure)
		agenda = agenda[1:]

		nextStates := make([]*State, 0)

		if nextState := state.OpenValve(); nextState != nil {
			nextStates = append(nextStates, nextState)
		}

		for _, next := range state.valve.leadsTo {
			nextState := state.MoveTo(next)
			nextStates = append(nextStates, nextState)
		}

		if state.time + 1 == endTime {
			for _, nextState := range nextStates {
				if nextState.pressure > max {
					max = nextState.pressure
					//fmt.Println("New max:", max)
				}
			}
		} else {
			for _, nextState := range nextStates {
				key := nextState.Key()
				bestPressure, found := best[key]
				if !found {
					//fmt.Printf("Pressure %d at %s not seen before\n", nextState.pressure, key, nextState.openValves);
				} else if bestPressure < nextState.pressure {
					//fmt.Printf("Pressure %d at %s/%x up from %d\n", nextState.pressure, key, bestPressure);
				} else {
					continue
				}
				agenda = append(agenda, nextState)
				best[key] = nextState.pressure
			}
		}
	}

	return max
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

	openValves := this.openValves
	openValves |= this.valve.bit
	time := this.time + 1
	pressure := this.pressure + this.rate
	rate := this.rate + this.valve.rate

	return &State{ this.valve, openValves, time, rate, pressure }
}

func (this *State) Key() string {
	return fmt.Sprintf("%s/%x", this.valve.name, this.openValves)
}
