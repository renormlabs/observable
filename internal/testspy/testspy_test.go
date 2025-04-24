// Copyright (c) 2025 Renorm Labs. All rights reserved.

package testspy_test

import (
	"testing"

	"renorm.dev/observable/internal/testspy"
)

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got nil")
		}
	}()
	f()
}

func TestSpySoftFail(t *testing.T) {
	spy := testspy.New(t)

	if spy.SpiedOnFailure {
		t.Fatalf("new SpyTB should start with Failed == false")
	}

	spy.Error("msg")
	if !spy.SpiedOnFailure {
		t.Fatalf("Error should set Failed flag")
	}

	// reset flag and test Errorf
	spy.SpiedOnFailure = false
	spy.Errorf("msg %d", 1)
	if !spy.SpiedOnFailure {
		t.Fatalf("Errorf should set Failed flag")
	}

	// reset and test Fail
	spy.SpiedOnFailure = false
	spy.Fail()
	if !spy.SpiedOnFailure {
		t.Fatalf("Fail should set Failed flag")
	}
}

func TestSpyHardFailPanics(t *testing.T) {
	spy := testspy.New(t)

	mustPanic(t, func() { spy.FailNow() })
	mustPanic(t, func() { spy.Fatal("boom") })
	mustPanic(t, func() { spy.Fatalf("boom %d", 1) })
}
