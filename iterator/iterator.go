package iterator

import (
    "bufio"
    "log"
    "os"
)

//------------------------------------------------------------------------------

type Stream[T any] interface {
    NilStream[T] | StaticStream[T] | DynamicStream[T]
}

type NilStream[T any] struct {

}

type StaticStream[T any, S Stream[T]] struct {
    value T
    next *S
}

type DynamicStream[T any] struct {
    value T
    promise func() *Stream[T]
}

//------------------------------------------------------------------------------

type Iterator[T any] interface {
    Next() bool
    Value() T
}

//------------------------------------------------------------------------------

type SliceIterator[T any] struct {
    slice []T
    value T
    index int
}

func MakeSliceIterator[T any](slice []T) Iterator[T]{
    return &SliceIterator[T]{
        slice: slice,
    }
}

func (iter *SliceIterator[T]) Value() T {
    return iter.value
}

func (iter *SliceIterator[T]) Next() bool {
    if (iter.index < len(iter.slice)) {
        iter.value = iter.slice[iter.index]
        iter.index++
        return true
    }
    return false
}

//------------------------------------------------------------------------------

type FilterIterator[T any] struct {
    source Iterator[T]
    predicate func (T) bool
}

func (iter *FilterIterator[T]) Value() T {
    return iter.source.Value()
}

func (iter *FilterIterator[T]) Next() bool {
    for iter.source.Next() {
        if iter.predicate(iter.source.Value()) {
            return true
        }
    }
    return false
}

func Filter[T any](source Iterator[T], predicate func(T) bool) Iterator[T] {
    return &FilterIterator[T]{
        source: source,
        predicate: predicate,
    }
}

//------------------------------------------------------------------------------

type MapIterator[S, T any] struct {
    source Iterator[S]
    transform func (S) T
    value T
}

func (iter *MapIterator[S, T]) Value() T {
    return iter.value
}

func (iter *MapIterator[S, T]) Next() bool {
    if (iter.source.Next()) {
        iter.value = iter.transform(iter.source.Value())
        return true
    }
    return false
}

func Map[S, T any](source Iterator[S], transform func(S) T) Iterator[T] {
    return &MapIterator[S, T]{
        source: source,
        transform: transform,
    }
}

//------------------------------------------------------------------------------

type TakeIterator[T any] struct {
    source Iterator[T]
    remaining int
}

func (iter *TakeIterator[T]) Value() T {
    return iter.source.Value()
}

func (iter *TakeIterator[T]) Next() bool {
    if (iter.remaining > 0) && iter.source.Next() {
        iter.remaining--
        return true
    }
    return false
}

func Take[T any](source Iterator[T], howMany int) Iterator[T] {
    return &TakeIterator[T]{
        source: source,
        remaining: howMany,
    }
}

//------------------------------------------------------------------------------

type ZipIterator[S, T, U any] struct {
    lhs Iterator[S]
    rhs Iterator[T]
    zip func(S, T) U
}

func (iter *ZipIterator[S, T, U]) Value() U {
    return iter.zip(iter.lhs.Value(), iter.rhs.Value())
}

func (iter *ZipIterator[S, T, U]) Next() bool {
    return iter.lhs.Next() && iter.rhs.Next()
}

func Zip[S, T, U any](lhs Iterator[S], rhs Iterator[T], zip func(S, T) U) Iterator[U] {
    return &ZipIterator[S, T, U]{
        lhs: lhs,
        rhs: rhs,
        zip: zip,
    }
}

//------------------------------------------------------------------------------

func Drop[T any](source Iterator[T], howMany int) Iterator[T] {
    for i := 0; i < howMany; i++ {
        source.Next()
    }
    return source
}

//------------------------------------------------------------------------------

func Foldl[S, T any](iter Iterator[S], f func(acc T, item S) T, base T) T {
    acc := base
    for (iter.Next()) {
        acc = f(acc, iter.Value())
    }
    return acc
}

//------------------------------------------------------------------------------

func Reduce[T any](iter Iterator[T], f func(acc T, item T) T) T {
    iter.Next()
    return Foldl(iter, f, iter.Value())
}

//------------------------------------------------------------------------------

func FileIterator(filename string) Iterator[string] {
	file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return MakeSliceIterator(lines)
}
