// Copyright (c) 2025 Renorm Labs. All rights reserved.

// Package observable provides lightweight, zero-dependency helpers for writing
// expressive assertions in Go tests. It embraces Go's standard testing package
// by integrating with [testing.TB] and keeps the API minimal while still
// supporting generic, user-defined predicates.
package observable

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// Predicate encapsulates a lazily‑evaluated boolean condition together with a descriptive failure message.
//
// Generally, users will construct predicates with the helper functions in this package (e.g. [Nil], [Equal], [Panics]).The zero value is NOT valid.
type Predicate struct {
	ok  func() bool
	msg func() string
}

// Ok evaluates and returns the underlying boolean condition.
func (p Predicate) Ok() bool { return p.ok() }

// Message returns the descriptive text explaining why the predicate failed.
func (p Predicate) Message() string { return p.msg() }

// Assertion is the constraint accepted by [Assert], [Assertf], and [Not]. It may be one of the following concrete types:
//   - bool          – a pre‑computed truth value
//   - func() bool – a zero‑arg function evaluated lazily
//   - [Predicate] – a value returned by this package's helpers
type Assertion interface {
	bool | func() bool | Predicate
}

// Assert evaluates the assertion and records an error on the [testing.TB] when the assertion is false.
//
// The returned bool is the evaluation result, which allows further composition or chaining inside a test when desired.
func Assert[T Assertion](tb testing.TB, a T) bool {
	tb.Helper()

	ok, msg := assert(a)

	return observe(tb, ok, msg) // auto-message (if any) is used
}

// Assertf behaves like [Assert] but lets the caller supply an explicit failure message via format and args, similar to [fmt.Sprintf].
func Assertf[T Assertion](tb testing.TB, a T, format string, args ...any) bool {
	tb.Helper()

	ok, _ := assert(a) // discard auto-message

	return observe(tb, ok, fmt.Sprintf(format, args...))
}

// Not returns the logical negation of its argument.
//
// The type bounds of the argument to Not are loose and enforced at runtime. You can:
//   - Negate a bool, resulting in a new bool
//   - Negate a func() bool, resulting in a new func() bool
//   - Negate a Predicate, resulting in a new predicate
//   - Negate a func(...) -> Assertion, resulting in new func(...) -> Assertion. This branch uses reflection and is slower.
//
// If you call Not with another type, it will panic!
func Not[F any](x F) F {
	//nolint:forcetypeassert
	switch v := any(x).(type) {
	case Predicate:
		return any(Predicate{
			ok:  func() bool { return !v.Ok() },
			msg: func() string { return "not: " + v.Message() },
		}).(F)
	case bool:
		return any(!v).(F)
	case func() bool:
		return any(func() bool { return !v() }).(F)
	}

	rv := reflect.ValueOf(x)
	if !rv.IsValid() || rv.Kind() != reflect.Func {
		panic(fmt.Sprintf("argument to Not must be an assertion or function that returns an assertion, got %v", rv.Kind()))
	}

	rt := rv.Type()
	if rt.NumOut() != 1 {
		panic(fmt.Sprintf("argument to Not must be an assertion or function that returns an assertion, got a function which returns %d items", rt.NumOut()))
	}

	out := rt.Out(0)
	switch out {
	case reflect.TypeOf(false):
		//nolint:forcetypeassert
		return reflect.MakeFunc(rt, func(args []reflect.Value) []reflect.Value {
			b := rv.Call(args)[0].Bool()
			return []reflect.Value{reflect.ValueOf(!b)}
		}).Interface().(F)
	case reflect.TypeOf((func() bool)(nil)):
		//nolint:forcetypeassert
		return reflect.MakeFunc(rt, func(args []reflect.Value) []reflect.Value {
			f := rv.Call(args)[0].Interface().(func() bool)
			return []reflect.Value{reflect.ValueOf(func() bool { return !f() })}
		}).Interface().(F)
	case reflect.TypeOf(Predicate{}):
		//nolint:forcetypeassert
		return reflect.MakeFunc(rt, func(args []reflect.Value) []reflect.Value {
			p := rv.Call(args)[0].Interface().(Predicate)
			neg := Predicate{
				ok:  func() bool { return !p.Ok() },
				msg: func() string { return "not: " + p.Message() },
			}
			return []reflect.Value{reflect.ValueOf(neg)}
		}).Interface().(F)
	default:
		panic(fmt.Sprintf("argument to Not must be an assertion or function that returns an assertion, got a function which returns %v", out))
	}
}

// Nil returns a [Predicate] that succeeds when v is nil.
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

// Zero returns a [Predicate] that succeeds when v is the zero value of its type.
func Zero[T comparable](v T) Predicate {
	return Predicate{
		ok:  func() bool { return v == *new(T) },
		msg: func() string { return fmt.Sprintf("expected zero value, got %v", v) },
	}
}

// Equal returns a [Predicate] that succeeds when got == want.
func Equal[T comparable](got, want T) Predicate {
	return Predicate{
		ok:  func() bool { return got == want },
		msg: func() string { return fmt.Sprintf("expected %v, got %v", want, got) },
	}
}

// Returns returns a [Predicate] that succeeds when f's return value equals want.
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

// observe is the common implementation used by [Assert] and [Assertf]. It reports a test error on tb when ok is false and returns ok so the caller can use the result in further logic.
//
//go:inline
func observe(tb testing.TB, ok bool, message string) bool {
	tb.Helper()

	if ok {
		return true
	}

	tb.Error(message)

	return false
}

// assert normalises any [Assertion] into its boolean value and accompanying auto‑generated message (if available).
func assert[T Assertion](a T) (ok bool, autoMsg string) {
	var zero T
	if _, isBool := any(zero).(bool); isBool {
		//nolint:forcetypeassert
		return any(a).(bool), ""
	}

	if _, isFunc := any(zero).(func() bool); isFunc {
		//nolint:forcetypeassert
		return any(a).(func() bool)(), ""
	}

	//nolint:forcetypeassert
	p := any(a).(Predicate)

	return p.Ok(), p.Message()
}
