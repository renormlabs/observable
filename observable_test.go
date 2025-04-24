// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable_test

import (
	"errors"
	"testing"

	"renorm.dev/observable"
	"renorm.dev/observable/internal/testspy"
)

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")
)

func expectPass[T observable.Assertion](tb testing.TB, pred T) {
	tb.Helper()
	spy := testspy.New(tb)
	if !observable.Assert(spy, pred) || spy.SpiedOnFailure {
		tb.Fatalf("expected pass, got fail")
	}
}

func expectFail[T observable.Assertion](tb testing.TB, pred T) {
	tb.Helper()
	spy := testspy.New(tb)
	if observable.Assert(spy, pred) || !spy.SpiedOnFailure {
		tb.Fatalf("expected fail, got pass")
	}
}

func TestBoolChecks(t *testing.T) {
	expectPass(t, true)
	expectFail(t, false)

	expectPass(t, observable.Not(false))
	expectFail(t, observable.Not(true))

	expectPass(t, func() bool { return 2+2 == 4 })
	expectFail(t, func() bool { return 2+2 == 5 })

	expectPass(t, observable.Not(func() bool { return 2+2 == 5 }))
	expectFail(t, observable.Not(func() bool { return 2+2 == 4 }))
}

func TestNilChecks(t *testing.T) {
	expectPass(t, observable.Nil(nil))
	expectFail(t, observable.Nil(1))

	expectPass(t, observable.Not(observable.Nil(1)))
	expectFail(t, observable.Not(observable.Nil(nil)))

	expectPass(t, observable.Not(observable.Nil)(1))
	expectFail(t, observable.Not(observable.Nil)(nil))
}

func TestErrorIsChecks(t *testing.T) {
	expectPass(t, observable.ErrorIs(errFoo, errFoo))
	expectFail(t, observable.ErrorIs(errFoo, errBar))

	expectPass(t, observable.Not(observable.ErrorIs)(errFoo, errBar))
	expectFail(t, observable.Not(observable.ErrorIs)(errFoo, errFoo))
}

func TestErrorsChecks(t *testing.T) {
	expectPass(t, observable.Errors(func() error { return errFoo }))
	expectFail(t, observable.Errors(func() error { return nil }))

	expectPass(t, observable.Not(observable.Errors)(func() error { return nil }))
	expectFail(t, observable.Not(observable.Errors)(func() error { return errFoo }))
}

func TestErrorsWithChecks(t *testing.T) {
	expectPass(t, observable.ErrorsWith(func() error { return errFoo }, errFoo))
	expectFail(t, observable.ErrorsWith(func() error { return errFoo }, errBar))

	expectPass(t, observable.Not(observable.ErrorsWith)(func() error { return errFoo }, errBar))
	expectFail(t, observable.Not(observable.ErrorsWith)(func() error { return errFoo }, errFoo))
}

func TestPanicsChecks(t *testing.T) {
	expectPass(t, observable.Panics(func() { panic("boom") }))
	expectFail(t, observable.Panics(func() {}))

	expectPass(t, observable.Not(observable.Panics)(func() {}))
	expectFail(t, observable.Not(observable.Panics)(func() { panic("boom") }))
}

func TestZeroChecks(t *testing.T) {
	expectPass(t, observable.Zero(""))
	expectFail(t, observable.Zero("foo"))

	expectPass(t, observable.Not(observable.Zero[string])("string"))
	expectFail(t, observable.Not(observable.Zero[string])(""))
}

func TestEqualChecks(t *testing.T) {
	expectPass(t, observable.Equal("a", "a"))
	expectFail(t, observable.Equal("a", "b"))

	expectPass(t, observable.Not(observable.Equal[string])("a", "b"))
	expectFail(t, observable.Not(observable.Equal[string])("a", "a"))
}

func TestReturnsChecks(t *testing.T) {
	// passing
	count := 0
	incr := func() int { count++; return 7 }
	expectPass(t, observable.Returns(incr, 7))
	if count != 1 {
		t.Fatalf("Returns should call function once, got %d", count)
	}

	// failing + not
	expectFail(t, observable.Returns(func() int { return 1 }, 2))

	expectPass(t, observable.Not(observable.Returns[int])(func() int { return 1 }, 2))
	expectFail(t, observable.Not(observable.Returns[int])(func() int { return 1 }, 1))
}

func TestAssertfOverride(t *testing.T) {
	spy := testspy.New(t)
	if observable.Assertf(spy, observable.Nil(1), "ignored") || !spy.SpiedOnFailure {
		t.Fatal("Assertf with failing predicate should fail")
	}

	spy = testspy.New(t)
	if !observable.Assertf(spy, observable.Nil(nil), "ignored") || spy.SpiedOnFailure {
		t.Fatal("Assertf with passing predicate should pass")
	}
}

func TestNotUnsupportedValue(t *testing.T) {
	got := observable.Not(42)
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestNotUnsupportedFunc(t *testing.T) {
	foo := func() int { return 7 }
	nilFunc := observable.Not(foo)

	if nilFunc != nil {
		t.Fatal("expected nil func, got non-nil")
	}
}
