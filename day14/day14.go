package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

type Material rune

const (
	Air    Material = '.'
	Rock   Material = '#'
	Sand   Material = 'o'
	Source Material = '+'
)

type Cave aoc.InfiniteGrid[Material]

type Vec2 struct {
	x, y int
}

func part1(lines []string) int {
	cave := parseCave(lines)

	start := Vec2{500, 0}
	cave.Set(start.x, start.y, Source)
	//cave.Print("%c", 0)
	sands := 0 // booyakasha
	for emitSand(cave, start) {
		sands++
		//cave.Print("%c", 0)
	}
	//cave.Print("%c", 0)
	return sands
}

func part2(lines []string) int {
	cave := parseCave(lines)

	sands := 0 // booyakasha

	start := Vec2{500, 0}
	cave.Set(start.x, start.y, Sand)
	sands++

	endY := cave.MaxY() + 2

	for y := 1; y < endY; y++ {
		for dx := -y; dx <= y; dx++ {
			x := dx + start.x
			if cave.Get(x, y) == Air && (cave.Get(x-1, y-1) == Sand ||
				cave.Get(x, y-1) == Sand ||
				cave.Get(x+1, y-1) == Sand) {
				sands++
				cave.Set(x, y, Sand)
			}
		}
	}

	return sands
}

func emitSand(cave *aoc.InfiniteGrid[Material], start Vec2) bool {
	pos := start
	for cave.OnGrid(pos.x, pos.y) {
		nextPos := moveSand(*cave, pos)
		if nextPos == pos {
			if cave.Get(pos.x, pos.y) == Sand {
				return false
			}
			cave.Set(pos.x, pos.y, Sand)
			return true
		}
		pos = nextPos
	}
	return false
}

func moveSand(cave aoc.InfiniteGrid[Material], from Vec2) Vec2 {
	if below := from.add(Vec2{0, 1}); cave.Get(below.x, below.y) == Air {
		return below
	}

	if sw := from.add(Vec2{-1, 1}); cave.Get(sw.x, sw.y) == Air {
		return sw
	}

	if se := from.add(Vec2{+1, 1}); cave.Get(se.x, se.y) == Air {
		return se
	}

	return from
}

func parseCave(lines []string) *aoc.InfiniteGrid[Material] {
	cave := aoc.NewInfiniteGrid[Material](Air)

	for _, line := range lines {
		vecs := parseLine(line)

		cave.Set(vecs[len(vecs)-1].x, vecs[len(vecs)-1].y, Rock)

		for pos, vecs := vecs[0], vecs[1:]; len(vecs) > 0; vecs = vecs[1:] {
			for delta := vecs[0].sub(pos).sign(); pos != vecs[0]; pos = pos.add(delta) {
				cave.Set(pos.x, pos.y, Rock)
			}
		}
	}

	return cave
}

func parseLine(line string) []Vec2 {
	pairs := strings.Split(line, " -> ")
	vecs := make([]Vec2, len(pairs))
	for i, pair := range pairs {
		vecs[i] = parseVec2(pair)
	}
	return vecs
}

func parseVec2(pair string) Vec2 {
	nums := aoc.ParseInts(strings.Split(pair, ","))
	return Vec2{nums[0], nums[1]}
}

func (this Vec2) add(that Vec2) Vec2 {
	return Vec2{this.x + that.x, this.y + that.y}
}

func (this Vec2) sub(that Vec2) Vec2 {
	return Vec2{this.x - that.x, this.y - that.y}
}

func (this Vec2) sign() Vec2 {
	return this.clamp(Vec2{-1, -1}, Vec2{1, 1})
}

func (this Vec2) clamp(min Vec2, max Vec2) Vec2 {
	return Vec2{clamp(this.x, min.x, max.x), clamp(this.y, min.y, max.y)}
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
