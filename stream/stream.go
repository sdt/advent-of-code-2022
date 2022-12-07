package stream

import (
    //"bufio"
    //"log"
    //"os"
)

//------------------------------------------------------------------------------

type Stream[T any] struct {
    value T
    next *Stream[T]
    promise func() *Stream[T]
}

func (this *Stream[T]) Value() T {
    return this.value
}

func (this *Stream[T]) Next() *Stream[T] {
    if this.next != nil {
        return this.next
    }

    if this.promise != nil {
        this.next = this.promise()
        this.promise = nil
    }

    return this.next
}

func NewStreamN[T any](value T, next *Stream[T]) *Stream[T] {
    return &Stream[T]{
        value: value,
        next: next,
    }
}

func NewStreamP[T any](value T, promise func() *Stream[T]) *Stream[T] {
    return &Stream[T]{
        value: value,
        promise: promise,
    }
}

//------------------------------------------------------------------------------

func FromSlice[T any](slice []T) *Stream[T] {
    if len(slice) == 0 {
        return nil
    }
    return NewStreamP[T](
        slice[0],
        func() *Stream[T] { return FromSlice(slice[1:]) },
    )
}

func Map[S, T any](s *Stream[S], transform func(S) T) *Stream[T] {
    if s == nil {
        return nil
    }
    return NewStreamP[T](
        transform(s.Value()),
        func() *Stream[T] { return Map[S, T](s.Next(), transform) },
    )
}

func Filter[T any](s *Stream[T], predicate func(T) bool) *Stream[T] {
    for s != nil {
        if predicate(s.Value()) {
            return NewStreamP[T](
                s.Value(),
                func() *Stream[T] { return Filter[T](s.Next(), predicate) },
            )
        }
        s = s.Next()
    }
    return s
}

func Concat[T any](ss []*Stream[T]) *Stream[T] {
    for {
        if len(ss) == 0 {
            return nil
        }

        if ss[0] == nil {
            ss = ss[1:]
        } else {
            break
        }
    }

    value := ss[0].value
    ss[0] = ss[0].Next()

    return NewStreamP[T](
        value,
        func() *Stream[T] { return Concat(ss) },
    )
}

func Append[T any](s *Stream[T], t *Stream[T]) *Stream[T] {
    if s == nil {
        return t
    }

    value := s.Value()
    s = s.Next()

    return NewStreamP[T](
        value,
        func() *Stream[T] { return Append(s, t) },
    );
}

func FlatMap[S, T any](s *Stream[S], transform func(S) []T) *Stream[T] {
    if s == nil {
        return nil
    }

    values := transform(s.Value())
    s = s.Next()

    if len(values) == 0 {
        return FlatMap(s, transform)
    }

    return Append(FromSlice(values), FlatMap(s, transform))
}

func Take[T any](s *Stream[T], howMany int) *Stream[T] {
    if s == nil {
        return nil
    }

    if howMany <= 0 {
        return nil
    }

    value := s.Value()
    s = s.Next()

    return NewStreamP[T](
        value,
        func() *Stream[T] { return Take(s, howMany-1) },
    );
}
