// Copyright (c) 2025 Renorm Labs. All rights reserved.

// Package testspy implements a lightweight wrapper around testing.TB to assist
// in the testing of testing frameworks.
package testspy

import (
	"testing"

	"renorm.dev/observable"
)

// SpyTB is a lightweight wrapper around testing.TB.
type SpyTB struct {
	testing.TB
	SpiedOnFailure bool
}

// New creates a new SpyTB instance from a testing.TB.
func New(t testing.TB) *SpyTB { return &SpyTB{TB: t} }

// Error intercepts calls to the regular Error method to mark test failure.
func (s *SpyTB) Error(...any) { s.SpiedOnFailure = true }

// Errorf intercepts calls to the regular Errorf method to mark test failure.
func (s *SpyTB) Errorf(string, ...any) { s.SpiedOnFailure = true }

// Fail intercepts calls to the regular Fail method to mark test failure.
func (s *SpyTB) Fail() { s.SpiedOnFailure = true }

// FailNow panics as this is not supported by SpyTB.
func (s *SpyTB) FailNow() { panic("FailNow not implemented on SpyTB") }

// Fatal panics as this is not supported by SpyTB.
func (s *SpyTB) Fatal(...any) { panic("Fatal not implemented on SpyTB") }

// Fatalf panics as this is not supported by SpyTB.
func (s *SpyTB) Fatalf(string, ...any) { panic("Fatalf not implemented on SpyTB") }

// ExpectPass expects an assertion to pass. Useful for testing a testing library.
func ExpectPass[T observable.Assertion](tb testing.TB, pred T) {
	tb.Helper()
	spy := New(tb)

	if !observable.Assert(spy, pred) || spy.SpiedOnFailure {
		switch x := any(pred).(type) {
		case observable.Predicate:
			tb.Errorf("expected pass, got fail: %v", x.Message())
		default:
			tb.Errorf("expected pass, got fail")
		}
	}
}

// ExpectFail expects an assertion to fail. Useful for testing a testing library.
func ExpectFail[T observable.Assertion](tb testing.TB, pred T) {
	tb.Helper()
	spy := New(tb)

	if observable.Assert(spy, pred) || !spy.SpiedOnFailure {
		tb.Errorf("expected fail, got pass")
	}
}
