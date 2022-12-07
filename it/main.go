package main

import (
    //"advent-of-code/it2"
    "advent-of-code/stream"
    "fmt"
)

func main() {
    s := stream.FromSlice([]int{ 1, 2, 3, 4, 5 })
    t := stream.Map(s, func(x int) int { return 3 * x })
    u := stream.Filter(t, func(x int) bool { return x & 1 == 1 })
    v := stream.Concat([]*stream.Stream[int]{ s, t, nil, s })

    printStream(s)
    printStream(t)
    printStream(u)
    printStream(v)

    printStream(stream.FlatMap(s, func (x int) []int {
        if x & 1 == 1 {
            return []int{ x }
        }
        return []int{ }
    }))

    printStream(stream.FlatMap(s, func (x int) []int {
        return []int{ x * 2 }
    }))

    xx := stream.FlatMap(s, func (x int) []int {
        if x & 1 == 1 {
            return []int{ }
        }
        ret := make([]int, x)
        for i := 0; i < x; i++ {
            ret[i] = x
        }
        return ret
    })
    printStream(xx)
    printStream(stream.Concat([]*stream.Stream[int]{ xx, xx }))
    printStream(stream.Take(stream.Append(xx, xx), 3))

    var integers *stream.Stream[int]
    integers = stream.NewStreamP[int](
        1,
        func() *stream.Stream[int] { return stream.Map(integers, func(x int) int { return x + 1 }) },
    )

    printStream(stream.Take(integers, 10))
}

func printStream[T any](s *stream.Stream[T]) {
    for s != nil {
        fmt.Println(s.Value())
        s = s.Next()
    }
    fmt.Println()
}

/*
func add(a, b int) int {
    return a + b
}

func main() {
    slice := []int{1, 2, 3, 4, 5}

    a := it2.SliceIt(slice)
    b := it2.FilterIt(a, func (n int) bool {
        return n & 1 == 1
    })
    c := it2.MapIt(b, func (n int) string {
        ret := ""
        for i := 0; i < n; i++ {
            ret += "."
        }
        return ret
    })

    d := it2.AccumIt(c, func (xs []string, x string) []string {
        return append(xs, x)
    }, []string{})

    e := it2.Take(d, 2)

    for _, value := range it2.Collect(e) {
        fmt.Println(value)
    }

    ints := it2.AccumIt(it2.Const(1), add, 0)

    for _, value := range it2.Collect(it2.Take(ints, 10)) {
        fmt.Println(value)
    }
}
*/
