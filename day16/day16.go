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
	valves []*Valve
	openValves uint64
	time int
	rate int
	pressure int
}

type StateSet map[string]*State

var seen StateSet

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)
	aa := parseInput(lines)

	fmt.Println(run(aa, 1, 30))
	fmt.Println(run(aa, 2, 26))
}

func run(start *Valve, actors int, endTime int) int {
	current := make(StateSet)
	seen = make(StateSet)

	AddState(current, NewState(start, actors))

	for time := 0; time < endTime-1; time++ {
		for _, state := range current {
			state.Tick()
		}

		for i := 0; i < actors; i++ {
			next := make(StateSet)
			for _, state := range current {
				AddState(next, state.OpenValve(i))

				for _, dest := range state.valves[i].leadsTo {
					AddState(next, state.MoveTo(i, dest))
				}
			}
			current = next
		}
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

	if existing, found := seen[key]; found {
		if existing.pressure >= state.pressure {
			return
		}
	}

	if existing, found := set[key]; found {
		if existing.pressure >= state.pressure {
			return
		}
	}
	set[key] = state
	seen[key] = state
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

func NewState(valve *Valve, actors int) *State {
	valves := make([]*Valve, actors)
	for i := 0; i < actors; i++ {
		valves[i] = valve
	}
	return &State{ valves, 0, 0, 0, 0 }
}

func (this *State) Tick() {
	this.time++
	this.pressure += this.rate
}

func (this *State) MoveTo(index int, dest *Valve) *State {
	valves := make([]*Valve, len(this.valves))
	copy(valves, this.valves)
	valves[index] = dest

	return &State{valves, this.openValves, this.time, this.rate, this.pressure}
}

func (this *State) OpenValve(index int) *State {
	valve := this.valves[index]
	if valve.rate == 0 {
		return nil
	}

	if this.openValves & valve.bit == valve.bit {
		return nil
	}

	openValves := this.openValves | valve.bit
	rate := this.rate + valve.rate

	return &State{ this.valves, openValves, this.time, rate, this.pressure }
}

func (this *State) Key() string {
	presentValves := uint64(0)
	for _, valve := range this.valves {
		presentValves |= valve.bit
	}
	return fmt.Sprintf("%x/%x", presentValves, this.openValves)
}
