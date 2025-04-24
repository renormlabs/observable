// Copyright (c) 2025 Renorm Labs. All rights reserved.

// Package testspy implements a lightweight wrapper around testing.TB to assist
// in the testing of testing frameworks.
package testspy

import "testing"

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
