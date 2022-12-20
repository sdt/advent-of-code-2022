package main

import (
	"advent-of-code/aoc"
	"fmt"
)

const Debug = true

type Node[T any] struct {
	prev, next *Node[T]
	value T
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	head := MakeEmpty[int]()
	cursor := head
	values := make([]*Node[int], len(lines))
	var zero *Node[int]
	for i, line := range lines {
		value := aoc.ParseInt(line)
		values[i] = &Node[int]{ value: value }
		cursor = values[i].InsertAfter(cursor)
		if value == 0 {
			zero = cursor
		}
	}
	head.Unlink()

	dump := func(start *Node[int]) {
		return
		start.Walk(func (node *Node[int]) {
			fmt.Printf("%d ", node.value)
		})
		fmt.Println()
	}

	dump(zero)

	for _, node := range values {
		node.Move(node.value)
		dump(zero)
	}

	total := 0
	for i := 0; i < 3; i++ {
		zero = zero.Next(1000)
		total += zero.value
	}

	return total
}

func part2(lines []string) int {
	key := 811_589_153
	head := MakeEmpty[int]()
	cursor := head
	values := make([]*Node[int], len(lines))
	var zero *Node[int]
	for i, line := range lines {
		value := aoc.ParseInt(line) * key
		values[i] = &Node[int]{ value: value }
		cursor = values[i].InsertAfter(cursor)
		if value == 0 {
			zero = cursor
		}
	}
	head.Unlink()

	dump := func(start *Node[int]) {
		return
		start.Walk(func (node *Node[int]) {
			fmt.Printf("%d ", node.value)
		})
		fmt.Println()
	}

	dump(zero)

	mod := len(values) - 1
	for i := 0; i < 10; i++ {
		for _, node := range values {
			if node.value > 0 {
				value := node.value % mod
				node.MoveRight(value)
			} else if node.value < 0 {
				value := -node.value % mod
				node.MoveLeft(value)
			}
		}
		dump(zero)
	}

	total := 0
	for i := 0; i < 3; i++ {
		zero = zero.Next(1000)
		total += zero.value
	}

	return total
}

func MakeEmpty[T any]() *Node[T] {
	var node Node[T]
	node.next = &node
	node.prev = &node
	return &node
}

func (this *Node[T]) Unlink() {
	this.next.prev = this.prev
	this.prev.next = this.next
}

func (this *Node[T]) Move(places int) {
	//fmt.Printf("Moving %v by %v\n", this, places)
	if places > 0 {
		this.MoveRight(places)
	} else if places < 0 {
		this.MoveLeft(-places)
	}
}

func (this *Node[T]) Next(places int) *Node[T] {
	for i := 0; i < places; i++ {
		this = this.next
	}
	return this
}

// Insert this after that
func (this *Node[T]) InsertAfter(that *Node[T]) *Node[T] {
	//fmt.Printf("Inserting %v after %v\n", this, that)
	this.prev = that
	this.next = that.next
	this.next.prev = this
	that.next = this
	return this
}

// Insert this before that
func (this *Node[T]) InsertBefore(that *Node[T]) *Node[T] {
	//fmt.Printf("Inserting %v before %v\n", this, that)
	this.next = that
	this.prev = that.prev
	that.prev.next = this
	that.prev = this
	return that
}

func (this *Node[T]) MoveRight(places int) {
	after := this.prev
	this.Unlink()
	for i := 0; i < places; i++ {
		after = after.next
	}
	this.InsertAfter(after)
}

func (this *Node[T]) MoveLeft(places int) {
	before := this.next
	this.Unlink()
	for i := 0; i < places; i++ {
		before = before.prev
	}
	this.InsertBefore(before)
}

func (this *Node[T]) Walk(visit func(node *Node[T])) {
	cursor := this

	for {
		visit(cursor)
		cursor = cursor.next
		if cursor == this {
			return
		}
	}
}
