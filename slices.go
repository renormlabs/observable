package observable

import (
	"fmt"
	"reflect"
	"sync"
)

// Length returns a [Predicate] that succeeds when len(s) == want for either a slice or a string.
func Length[T any](s []T, want int) Predicate {
	return Predicate{
		ok:  func() bool { return len(s) == want },
		msg: func() string { return fmt.Sprintf("expected length %d, got %d", want, len(s)) },
	}
}

// Empty returns a [Predicate] that succeeds when len(s) == 0 for slice.
func Empty[T any](s []T) Predicate { return Length(s, 0) }

// Contains returns a [Predicate] that succeeds when elem is present in slice.
func Contains[T comparable](slice []T, elem T) Predicate {
	return Predicate{
		ok: func() bool {
			for _, v := range slice {
				if v == elem {
					return true
				}
			}
			return false
		},
		msg: func() string { return fmt.Sprintf("expected %v to contain %v", slice, elem) },
	}
}

// SequenceEqual returns a [Predicate] that succeeds when got and want have identical length and elements appear in the same order.
func SequenceEqual[T comparable](got, want []T) Predicate {
	return Predicate{
		ok: func() bool {
			if len(got) != len(want) {
				return false
			}
			for i, v := range got {
				if v != want[i] {
					return false
				}
			}
			return true
		},
		msg: func() string { return fmt.Sprintf("expected slice %v, got %v", want, got) },
	}
}

// SequenceDeepEqual returns a [Predicate] that succeeds when want and got [reflect.DeepEqual] each other. Allows comparing slices with non-comparable element types.
func SequenceDeepEqual[T any](got, want []T) Predicate {
	var (
		once  sync.Once
		match bool
	)

	check := func() { once.Do(func() { match = reflect.DeepEqual(got, want) }) }

	return Predicate{
		ok: func() bool {
			check()
			return match
		},
		msg: func() string {
			check()
			return fmt.Sprintf("expected slice %v, got %v", want, got)
		},
	}
}

// ElementsMatch returns a [Predicate] that succeeds when the two slices contain the same multiset of elements, irrespective of order.
func ElementsMatch[T comparable](got, want []T) Predicate {
	count := func(s []T) map[T]int {
		m := make(map[T]int, len(s))
		for _, v := range s {
			m[v]++
		}

		return m
	}

	var (
		once  sync.Once
		match bool
	)

	check := func() { once.Do(func() { match = reflect.DeepEqual(count(got), count(want)) }) }

	return Predicate{
		ok: func() bool {
			check()
			return match
		},
		msg: func() string {
			check()
			return fmt.Sprintf("expected %v and %v to contain the same elements", got, want)
		},
	}
}
