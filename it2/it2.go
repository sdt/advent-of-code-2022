package it2

import (
    //"bufio"
    //"log"
    //"os"
)

//------------------------------------------------------------------------------

type Iterator[T any] func() (T, bool)

//------------------------------------------------------------------------------

func done[T any]() (T, bool) {
    var nothing T
    return nothing, false
}

func Slice[T any](slice []T) Iterator[T] {
    index := 0
    return func() (T, bool) {
        if index == len(slice) {
            return done[T]()
        }
        index++
        return slice[index-1], true
    }
}

func Map[S, T any](iter Iterator[S], transform func(S) T) Iterator[T] {
    return func() (T, bool) {
        value, valid := iter()
        if valid {
            return transform(value), valid
        }
        return done[T]()
    }
}

func Filter[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
    return func() (T, bool) {
        for value, valid := iter(); valid; value, valid = iter() {
            if predicate(value) {
                return value, true
            }
        }
        return done[T]()
    }
}

func Accumulate[S, T any](iter Iterator[S], f func(T, S) T, base T)  Iterator[T] {
    acc := base
    return func() (T, bool) {
        if value, valid := iter(); valid {
            acc = f(acc, value)
            return acc, true
        }
        return done[T]()
    }
}

func Take[T any](iter Iterator[T], howMany int) Iterator[T] {
    return func() (T, bool) {
        if howMany > 0 {
            howMany--
            if value, valid := iter(); valid {
                return value, true
            }
        }
        return done[T]()
    }
}

func TakeWhile[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
    taken := false
    return func() (T, bool) {
        if !taken {
            if value, valid := iter(); valid {
                if predicate(value) {
                    return value, true
                }
            }
            taken = true
        }
        return done[T]()
    }
}

func DropWhile[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
    dropped := false
    return func() (T, bool) {
        if dropped {
            return iter()
        }

        for {
            value, valid := iter()
            if !valid {
                return done[T]()
            }
            if !predicate(value) {
                dropped = true
                return value, true
            }
        }
    }
}

func Zip[S, T, U any](iterA Iterator[S], iterB Iterator[T], zip func(S, T) U) Iterator[U] {
    return func() (U, bool) {
        if valueA, validA := iterA(); validA {
            if valueB, validB := iterB(); validB {
                return zip(valueA, valueB), true
            }
        }
        return done[U]()
    }
}

func Interleave[T any](iters []Iterator[T]) Iterator[T] {
    index := 0
    return func() (T, bool) {
        if value, valid := iters[index](); valid {
            index = (index + 1) % len(iters)
            return value, valid
        }
        return done[T]()
    }
}

func Const[T any](value T) Iterator[T] {
    return func() (T, bool) {
        return value, true
    }
}

//------------------------------------------------------------------------------

func FoldL[S, T any](iter Iterator[S], f func(T, S) T, base T) T {
    acc := base
    for value, valid := iter(); valid; value, valid = iter() {
        acc = f(acc, value)
    }
    return acc
}

func Reduce[T any](iter Iterator[T], f func(T, T) T) T {
    if value, valid := iter(); valid {
        return FoldL(iter, f, value)
    }
    var none T
    return none
}

func Collect[T any](iter Iterator[T]) []T {
    return FoldL(iter, func(xs []T, x T) []T {
        return append(xs, x)
    }, []T{});
}
