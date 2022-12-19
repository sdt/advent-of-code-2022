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
	cost [MaterialCount]uint16
}

type Blueprint struct {
	robot [MaterialCount]RobotCost
}

type State struct {
	materialCount [MaterialCount]uint16
	robotCount [MaterialCount]uint16
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	blueprints := make([]Blueprint, len(lines))
	for i, line := range lines {
		blueprints[i] = parseBlueprint(line)
	}

	fmt.Println(part1(blueprints))
	fmt.Println(part2(blueprints[0:3]))
}

func part1(blueprints []Blueprint) int {
	totalQuality := 0
	for i, blueprint := range blueprints {
		quality := (i + 1) * runBlueprint(blueprint, 24)
		totalQuality += quality
	}
	return totalQuality
}

func part2(blueprints []Blueprint) int {
	result := 1
	for _, blueprint := range blueprints {
		result *= runBlueprint(blueprint, 32)
	}
	return result
}

// 3096 is too low

func runBlueprint(blueprint Blueprint, minutes int) int {
	currentStates := make(map[State]bool)
	currentStates[startState()] = true
	nextStates := make(map[State]bool)

	for minute := 1; minute <= minutes; minute++ {
		//fmt.Printf("Minute %d: states %d\n", minute, len(currentStates))
		for state := range currentStates {
			if state.canBuild(blueprint, Geode, 1) {
				nextStates[ state.build(blueprint, Geode) ] = true
				continue
			}

			if state.canBuild(blueprint, Obsidian, 1) {
				nextStates[ state.build(blueprint, Obsidian) ] = true
				continue
			}
			for material := Material(0); material <= Clay; material++ {
				// Can't build a robot if you can't afford it.
				// And if you can afford to build two, it's too late.
				if !state.canBuild(blueprint, material, 1) || state.canBuild(blueprint, material, 2) {
					continue
				}
				nextStates[ state.build(blueprint, material) ] = true
			}
			nextStates[state.step()] = true
		}
		currentStates = nextStates
		nextStates = make(map[State]bool)
	}

	max := uint16(0)
	for state := range currentStates {
		if state.materialCount[Geode] > max {
			max = state.materialCount[Geode]
		}
	}
	//fmt.Println("Max =", max)
	return int(max)
}

func parseBlueprint(line string) Blueprint {
	matches := InputRegex.FindAllString(line, -1)
	numbers := aoc.ParseInts(matches)
	var bp Blueprint

	// [RobotType][Material]
	bp.robot[Ore].cost[Ore]  = uint16(numbers[1])
	bp.robot[Clay].cost[Ore] = uint16(numbers[2])
	bp.robot[Obsidian].cost[Ore] = uint16(numbers[3])
	bp.robot[Obsidian].cost[Clay] = uint16(numbers[4])
	bp.robot[Geode].cost[Ore] = uint16(numbers[5])
	bp.robot[Geode].cost[Obsidian] = uint16(numbers[6])
	return bp
}

func startState() State {
	var state State
	state.robotCount[Ore] = 1
	return state
}

func (this State) canBuild(bp Blueprint, robot Material, count uint16) bool {
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
	for material, count := range this.robotCount {
		this.materialCount[material] += count
	}
	this.robotCount[robot]++
	return this
}

func (this State) step() State {
	for material, count := range this.robotCount {
		this.materialCount[material] += count
	}
	return this
}
