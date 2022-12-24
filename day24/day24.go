package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Wind uint8
const (
	North Wind = 1 << iota
	East
	South
	West
	None Wind = 0
	Blocked = 0xff
)

type Map struct {
	aoc.Grid[Wind]
}

type Vec2 struct {
	x, y int8
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	current := parseMap(lines)
	//current.Print()

	locations := make(map[Vec2]bool)
	locations[ current.startPos() ] = true
	minute := 0

	moves := [...]Vec2{ {0,0}, {-1,0}, {1,0}, {0,1}, {0,-1} }

	for len(locations) > 0 {
		minute++

		next := step(current)
		//next.Print()
		nextLocations := make(map[Vec2]bool)

		for location := range locations {
			for _, move := range moves {
				nextLocation := location.Add(move)
				if next.isGoal(nextLocation) {
					return minute
				}
				if next.canMove(nextLocation) {
					//fmt.Printf("Moving from %v to %v\n", location, nextLocation)
					nextLocations[nextLocation] = true
				}
			}
		}

		current = next
		locations = nextLocations

		//fmt.Printf("Minute %d: %d locations\n", minute, len(locations))
	}
	return minute
}

func makeMap(w, h int) *Map {
	return &Map{ *aoc.NewGrid[Wind](w, h, None) }
}

func parseMap(lines []string) *Map {
	w := len(lines[0]) - 2
	h := len(lines) - 2

	m := makeMap(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			wind := None
			switch lines[y+1][x+1] {
			case '>': wind |= East
			case '<': wind |= West
			case '^': wind |= North
			case 'v': wind |= South
			}
			m.Set(x, y, wind)
		}
	}
	return m
}

func (this *Map) AddWind(x, y int, wind Wind) {
	this.Set(x, y, wind | this.Get(x, y))
}

func step(current *Map) *Map {
	w := current.Width()
	h := current.Height()
	next := makeMap(current.Width(), current.Height())

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			wind := current.Get(x, y)

			if wind & North == North {
				next.AddWind(x, (y + h - 1) % h, North)
			}
			if wind & South == South {
				next.AddWind(x, (y + h + 1) % h, South)
			}
			if wind & West == West {
				next.AddWind((x + w - 1) % w, y, West)
			}
			if wind & East == East {
				next.AddWind((x + w + 1) % w, y, East)
			}
		}
	}

	return next
}

func (this *Map) Print() {
	w := this.Width()
	h := this.Height()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			wind := this.Get(x, y)
			var char rune

			switch wind {
			case North: char = '^'
			case South: char = 'V'
			case East:  char = '>'
			case West:  char = '<'
			case None:  char = '.'
			default: char = rune('0' + countBits(wind))
			}

			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (this *Map) isGoal(p Vec2) bool {
	return int(p.x) == this.Width() - 1 && int(p.y) == this.Height()
}

func (this *Map) startPos() Vec2 {
	return Vec2{ 0, -1 }
}

func (this *Map) canMove(p Vec2) bool {
	if p == this.startPos() {
		return true // home position
	}

	x := int(p.x)
	y := int(p.y)

	if !this.Contains(x, y) {
		return false
	}
	return this.Get(int(p.x), int(p.y)) == None
}

func countBits(w Wind) int {
	count := 0
	for i := 0; i < 4; i++ {
		if w & (1 << i) != 0 {
			count++
		}
	}
	return count
}

func (this Vec2) Add(that Vec2) Vec2 {
	return Vec2{ this.x + that.x, this.y + that.y }
}
