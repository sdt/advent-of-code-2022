package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

type Vec2 struct {
	x, y int
}

type Pair struct {
	sensor, beacon Vec2
	distance       int
}

type Segment struct {
	start, length int
}

func part1(lines []string) int {
	pairs := make([]*Pair, len(lines))

	for i, line := range lines {
		pairs[i] = parsePair(line)
	}

	line := 10
	//line := 2000000

	segments := GetSegments(pairs, line)
	beacons := make(map[int]bool)
	for _, pair := range pairs {
		if pair.beacon.y == line {
			beacons[pair.beacon.x] = true
		}
	}

	total := -len(beacons)
	for _, segment := range segments {
		total += segment.length
	}

	return total
}

func part2(lines []string) int {
	pairs := make([]*Pair, len(lines))

	for i, line := range lines {
		pairs[i] = parsePair(line)
		//fmt.Println(line, pairs[i])
	}

	//maxY := 20
	maxY := 4_000_000

	for y := maxY; y >= 0; y-- {
		segments := GetSegments(pairs, y)
		if len(segments) > 1 {
			x := segments[1].start - 1
			return y + x*4_000_000
		}
	}

	return 0
}

func GetSegments(pairs []*Pair, line int) []Segment {
	segments := []Segment{}
	for _, pair := range pairs {
		segment := pair.GetSegment(line)
		if segment.length > 0 {
			segments = segment.Insert(segments)
		}
	}
	return segments
}

func (this *Pair) GetSegment(y int) Segment {
	// If we have y == this.sensor.y we are at max width
	maxWidth := 2*this.distance + 1

	// Each row up or down from this.sensor.y reduces the width by 2
	width := maxWidth - 2*abs(y-this.sensor.y)

	// Too far up or down returns an empty segment
	if width < 0 {
		return Segment{}
	}

	// We don't care about the beacon for now.
	segment := Segment{this.sensor.x - width/2, width}
	return segment
}

func (this Segment) Insert(in []Segment) []Segment {
	out := make([]Segment, 0, len(in)+1)

	for len(in) > 0 {
		comp := this.Compare(&in[0])

		if comp < 0 {
			// this is entirely before in[0]
			out = append(out, this)
			out = append(out, in...)
			return out
		}
		if comp > 0 {
			// in[0] is entirely before this
			out = append(out, in[0])
		} else {
			// this overlaps in[0]
			this = this.MergeWith(in[0])
		}
		in = in[1:]
	}
	out = append(out, this)
	return out
}

func (this Segment) String() string {
	return fmt.Sprintf("<%d..%d>(%d)", this.start, this.End(), this.length)
}

func (this *Segment) Compare(that *Segment) int {
	if this.End()+1 < that.start {
		return -1
	}

	if that.End()+1 < this.start {
		return +1
	}

	return 0
}

func (this Segment) MergeWith(that Segment) Segment {
	if this.start > that.start {
		return that.MergeWith(this)
	}

	end := this.End()
	if that.End() > end {
		end = that.End()
	}

	return Segment{this.start, end - this.start + 1}
}

func (this *Segment) Length() int {
	return this.length
}

func (this *Segment) End() int {
	return this.start + this.length - 1
}

func manhattanDistance(a, b Vec2) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var InputRegex = regexp.MustCompile("Sensor at x=(-?\\d+), y=(-?\\d+): closest beacon is at x=(-?\\d+), y=(-?\\d+)")

func parsePair(line string) *Pair {
	matches := InputRegex.FindStringSubmatch(line)
	nums := aoc.ParseInts(matches[1:])

	sensor := Vec2{nums[0], nums[1]}
	beacon := Vec2{nums[2], nums[3]}

	return &Pair{
		sensor:   sensor,
		beacon:   beacon,
		distance: manhattanDistance(sensor, beacon),
	}
}
