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

// Assert evaluates the predicate and records an error on the [testing.TB] when the predicate is false.
//
// The returned bool is the evaluation result, which allows further composition or chaining inside a test when desired.
func Assert(tb testing.TB, p Predicate) bool {
	tb.Helper()

	return observe(tb, p.Ok(), p.Message())
}

// Assertf behaves like [Assert] but lets the caller supply an explicit failure message via format and args, similar to [fmt.Sprintf].
func Assertf(tb testing.TB, p Predicate, format string, args ...any) bool {
	tb.Helper()
	return observe(tb, p.Ok(), fmt.Sprintf(format, args...))
}

// That promotes a bool or bool-thunk to a [Predicate].
func That[T ~bool | ~func() bool](x T) Predicate {
	var (
		once sync.Once
		got  bool
	)

	eval := func() {
		once.Do(func() {
			if f, ok := any(x).(func() bool); ok {
				got = f()
			} else {
				got = any(x).(bool)
			}
		})
	}

	return Predicate{
		ok:  func() bool { eval(); return got },
		msg: func() string { eval(); return fmt.Sprintf("expected true, got %v", got) },
	}
}

// Not returns the logical negation of its argument.
//
// You can negate a:
// - [Predicate], resulting in a [Predicate]
// - A function of any arity that returns a Predicate
//
// Negating a function of positive arity will use runtime reflection.
//
// Calling Not with anything else will **panic**!
func Not[T any](a T) T {
	if p, ok := any(a).(Predicate); ok {
		return any(Predicate{
			ok:  func() bool { return !p.Ok() },
			msg: func() string { return fmt.Sprintf("not: %s", p.Message()) },
		}).(T)
	}

	// Handle nullary functions without reflection.
	if f, ok := any(a).(func() Predicate); ok {
		return any(func() Predicate {
			return Not(f())
		}).(T)
	}

	// Handle postive-arity functions with reflection.
	rv := reflect.ValueOf(a)
	if rv.Kind() != reflect.Func {
		panic(fmt.Sprintf("argument to Not must be Predicate or func→Predicate, got %v", rv.Kind()))
	}
	rt := rv.Type()
	if rt.NumOut() != 1 || rt.Out(0) != reflect.TypeOf(Predicate{}) {
		panic("argument to Not must return a Predicate")
	}

	wrapper := reflect.MakeFunc(rt, func(args []reflect.Value) []reflect.Value {
		var callArgs []reflect.Value
		if rt.IsVariadic() {
			numFixedArgs := rt.NumIn() - 1
			callArgs = make([]reflect.Value, 0, len(args)-1+5) // Pre-allocate: num fixed + estimate for variadic

			// Add fixed arguments
			for i := 0; i < numFixedArgs; i++ {
				callArgs = append(callArgs, args[i])
			}

			// Unpack and add variadic arguments
			// args[numFixedArgs] is the reflect.Value representing the slice of variadic arguments
			if len(args) > numFixedArgs { // Ensure the variadic slice argument itself is present
				variadicSliceValue := args[numFixedArgs]
				for i := 0; i < variadicSliceValue.Len(); i++ {
					callArgs = append(callArgs, variadicSliceValue.Index(i))
				}
			}
		} else {
			callArgs = args
		}

		out := rv.Call(callArgs)
		p := out[0].Interface().(Predicate) // Original function returned a Predicate

		negatedP := Predicate{
			ok:  func() bool { return !p.Ok() },
			msg: func() string { return fmt.Sprintf("not: %s", p.Message()) },
		}
		return []reflect.Value{reflect.ValueOf(negatedP)}
	})

	return wrapper.Interface().(T)
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
