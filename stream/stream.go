package stream

import (
// "fmt"
)

//------------------------------------------------------------------------------

type Promise[T any] func() *Stream[T]

type Next interface{}

type Stream[T any] struct {
	value T
	tail  Next // either *Stream[T] or Promise[T]
}

func (this *Stream[T]) Value() T {
	return this.value
}

func (this *Stream[T]) Next() *Stream[T] {
	// If the tail is a promise, call it and store the result
	if promise, ok := this.tail.(Promise[T]); ok {
		this.tail = promise()
	}

	return this.tail.(*Stream[T])
}

func (this *Stream[T]) NextValue() (T, *Stream[T]) {
	value := this.value

	return value, this.Next()
}

func Cons[T any](value T, promise Promise[T]) *Stream[T] {
	return &Stream[T]{
		value: value,
		tail:  promise,
	}
}

//------------------------------------------------------------------------------

func FromSlice[T any](slice []T) *Stream[T] {
	if len(slice) == 0 {
		return nil
	}
	return Cons[T](
		slice[0],
		func() *Stream[T] { return FromSlice(slice[1:]) },
	)
}

func Map[S, T any](s *Stream[S], transform func(S) T) *Stream[T] {
	if s == nil {
		return nil
	}
	return Cons[T](
		transform(s.Value()),
		func() *Stream[T] { return Map[S, T](s.Next(), transform) },
	)
}

func Filter[T any](s *Stream[T], predicate func(T) bool) *Stream[T] {
	for s != nil {
		if predicate(s.Value()) {
			return Cons[T](
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

	return Cons[T](
		value,
		func() *Stream[T] { return Concat(ss) },
	)
}

func Append[T any](s *Stream[T], t *Stream[T]) *Stream[T] {
	if s == nil {
		return t
	}

	value, s := s.NextValue()

	return Cons[T](
		value,
		func() *Stream[T] { return Append(s, t) },
	)
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

	value, s := s.NextValue()

	return Cons[T](
		value,
		func() *Stream[T] { return Take(s, howMany-1) },
	)
}

func Zip[S, T, U any](s *Stream[S], t *Stream[T], transform func(S, T) U) *Stream[U] {
	if s == nil || t == nil {
		return nil
	}

	sValue, s := s.NextValue()
	tValue, t := t.NextValue()

	return Cons[U](
		transform(sValue, tValue),
		func() *Stream[U] { return Zip(s, t, transform) },
	)
}

func Sieve(s *Stream[int]) *Stream[int] {
	value, s := s.NextValue()

	return Cons[int](
		value,
		func() *Stream[int] {
			return Sieve(Filter[int](s, func(x int) bool {
				return x%value != 0
			}))
		},
	)
}
