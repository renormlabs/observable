// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable

import (
	"fmt"
	"reflect"
	"sync"
)

// Nil returns a [Predicate] that is ok when v is nil.
func Nil(v any) Predicate {
	isNil := func(x any) bool {
		if x == nil {
			return true
		}

		rv := reflect.ValueOf(x)

		return rv.Kind() >= reflect.Chan && rv.Kind() <= reflect.Slice && rv.IsNil()
	}

	return Predicate{
		ok:  func() bool { return isNil(v) },
		msg: func() string { return fmt.Sprintf("expected %#v to be nil", v) },
	}
}

// Zero returns a [Predicate] that is ok when v is the zero value of its type.
func Zero[T comparable](v T) Predicate {
	return Predicate{
		ok:  func() bool { return v == *new(T) },
		msg: func() string { return fmt.Sprintf("expected zero value, got %v", v) },
	}
}

// Equal returns a [Predicate] that is ok when got == want.
func Equal[T comparable](got, want T) Predicate {
	return Predicate{
		ok:  func() bool { return got == want },
		msg: func() string { return fmt.Sprintf("expected %v, got %v", want, got) },
	}
}

// Returns returns a [Predicate] that is ok when f's return value equals want.
func Returns[T comparable](f func() T, want T) Predicate {
	var (
		once sync.Once
		got  T
	)

	eval := func() { once.Do(func() { got = f() }) }

	return Predicate{
		ok:  func() bool { eval(); return got == want },
		msg: func() string { eval(); return fmt.Sprintf("expected %v, got %v", want, got) },
	}
}

// True returns a Predicate that always is ok.
func True() Predicate {
	return Predicate{
		ok:  func() bool { return true },
		msg: func() string { return "true" },
	}
}

// False returns a Predicate that always is not ok.
func False() Predicate {
	return Predicate{
		ok:  func() bool { return false },
		msg: func() string { return "false" },
	}
}

// Any returns a [Predicate] that is ok when any of the supplied predicates are ok.
func Any(ps ...Predicate) Predicate {
	var (
		once sync.Once
		msgs []string
	)

	eval := func() {
		once.Do(func() {
			for _, p := range ps {
				if !p.Ok() {
					msgs = append(msgs, p.Message())
				}
			}
		})
	}

	return Predicate{
		ok:  func() bool { eval(); return len(msgs) < len(ps) },
		msg: func() string { eval(); return fmt.Sprintf("expected any to be true, all failed: %v", msgs) },
	}
}

// All returns a [Predicate] that is ok when all of the supplied predicates are ok.
func All(ps ...Predicate) Predicate {
	var (
		once sync.Once
		msgs []string
	)

	eval := func() {
		once.Do(func() {
			for _, p := range ps {
				if !p.Ok() {
					msgs = append(msgs, p.Message())
				}
			}
		})
	}

	return Predicate{
		ok:  func() bool { eval(); return len(msgs) == 0 },
		msg: func() string { eval(); return fmt.Sprintf("expected all to be true, failures: %v", msgs) },
	}
}
