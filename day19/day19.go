package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
)

var InputRegex = regexp.MustCompile("\\d+")

type Material = int

const (
	Ore Material = iota
	Clay
	Obsidian
	Geode

	MaterialCount
)

type RobotCost struct {
	cost [MaterialCount]int
}

type Blueprint struct {
	robot [MaterialCount]RobotCost
}

type State struct {
	materialCount [MaterialCount]int
	robotCount [MaterialCount]int
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	blueprints := make([]Blueprint, len(lines))
	for i, line := range lines {
		blueprints[i] = parseBlueprint(line)
	}

	fmt.Println(part1(blueprints))
}

func part1(blueprints []Blueprint) int {
	totalQuality := 0
	for i, blueprint := range blueprints {
		quality := (i + 1) * runBlueprint(blueprint)
		totalQuality += quality
	}
	return totalQuality
}

func runBlueprint(blueprint Blueprint) int {
	currentStates := make(map[State]bool)
	currentStates[startState()] = true
	nextStates := make(map[State]bool)

	for minute := 1; minute <= 23; minute++ {
		//fmt.Printf("Minute %d: states %d\n", minute, len(currentStates))
		for state := range currentStates {
			for material := Material(0); material < MaterialCount; material++ {
				// Can't build a robot if you can't afford it.
				// And if you can afford to build two, it's too late.
				if !state.canBuild(blueprint, material, 1) || state.canBuild(blueprint, material, 2) {
					continue
				}
				nextState := state.build(blueprint, material)
				nextState = nextState.step()
				nextState.robotCount[material]++

				nextStates[nextState] = true
			}
			nextStates[state.step()] = true
		}
		currentStates = nextStates
		nextStates = make(map[State]bool)
	}

	max := 0
	for state := range currentStates {
		state = state.step()
		if state.materialCount[Geode] > max {
			max = state.materialCount[Geode]
		}
	}
	//fmt.Println("Max =", max)
	return max
}

func parseBlueprint(line string) Blueprint {
	matches := InputRegex.FindAllString(line, -1)
	numbers := aoc.ParseInts(matches)
	var bp Blueprint

	// [RobotType][Material]
	bp.robot[Ore].cost[Ore]  = numbers[1]
	bp.robot[Clay].cost[Ore] = numbers[2]
	bp.robot[Obsidian].cost[Ore] = numbers[3]
	bp.robot[Obsidian].cost[Clay] = numbers[4]
	bp.robot[Geode].cost[Ore] = numbers[5]
	bp.robot[Geode].cost[Obsidian] = numbers[6]
	return bp
}

func startState() State {
	var state State
	state.robotCount[Ore] = 1
	return state
}

func (this State) canBuild(bp Blueprint, robot Material, count int) bool {
	for material, cost := range bp.robot[robot].cost {
		if cost * count > this.materialCount[material] {
			return false
		}
	}
	return true
}

func (this State) build(bp Blueprint, robot Material) State {
	for material, cost := range bp.robot[robot].cost {
		this.materialCount[material] -= cost
	}
	return this
}

func (this State) step() State {
	for material, count := range this.robotCount {
		this.materialCount[material] += count
	}
	return this
}
