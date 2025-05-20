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

func TestBoolChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.That(true))
	testspy.ExpectFail(t, observable.That(false))

	testspy.ExpectPass(t, observable.Not(observable.That(false)))
	testspy.ExpectFail(t, observable.Not(observable.That(true)))

	testspy.ExpectPass(t, observable.That(func() bool { return 2+2 == 4 }))
	testspy.ExpectFail(t, observable.That(func() bool { return 2+2 == 5 }))

	testspy.ExpectPass(t, observable.Not(observable.That(func() bool { return 2+2 == 5 })))
	testspy.ExpectFail(t, observable.Not(observable.That(func() bool { return 2+2 == 4 })))
}

func TestNilChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.Nil(nil))
	testspy.ExpectFail(t, observable.Nil(1))

	testspy.ExpectPass(t, observable.Not(observable.Nil(1)))
	testspy.ExpectFail(t, observable.Not(observable.Nil(nil)))

	testspy.ExpectPass(t, observable.Not(observable.Nil)(1))
	testspy.ExpectFail(t, observable.Not(observable.Nil)(nil))
}

func TestZeroChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.Zero(""))
	testspy.ExpectFail(t, observable.Zero("foo"))

	testspy.ExpectPass(t, observable.Not(observable.Zero[string])("string"))
	testspy.ExpectFail(t, observable.Not(observable.Zero[string])(""))
}

func TestEqualChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.Equal("a", "a"))
	testspy.ExpectFail(t, observable.Equal("a", "b"))

	testspy.ExpectPass(t, observable.Not(observable.Equal[string])("a", "b"))
	testspy.ExpectFail(t, observable.Not(observable.Equal[string])("a", "a"))
}

func TestReturnsChecks(t *testing.T) {
	// passing
	count := 0
	incr := func() int { count++; return 7 }
	testspy.ExpectPass(t, observable.Returns(incr, 7))

	if count != 1 {
		t.Fatalf("Returns should call function once, got %d", count)
	}

	// failing + not
	testspy.ExpectFail(t, observable.Returns(func() int { return 1 }, 2))

	testspy.ExpectPass(t, observable.Not(observable.Returns[int])(func() int { return 1 }, 2))
	testspy.ExpectFail(t, observable.Not(observable.Returns[int])(func() int { return 1 }, 1))
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
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic, did not panic")
		}
	}()

	observable.Not(42)
}

func TestNotUnsupportedFunc2(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic, did not panic")
		}
	}()

	foo := func() (observable.Predicate, error) { return observable.Equal("a", "b"), nil }
	observable.Not(foo)
}

func TestNotChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.Not(observable.Equal[string])("a", "b"))
	testspy.ExpectPass(t, observable.Not(observable.That(false)))
	testspy.ExpectPass(t, observable.Not(observable.That(func() bool { return false })))

	testspy.ExpectPass(t, observable.True())
	testspy.ExpectFail(t, observable.False())
	testspy.ExpectFail(t, observable.Not(observable.True)())
	testspy.ExpectPass(t, observable.Not(observable.False)())
}

func TestAny(t *testing.T) {
	testspy.ExpectPass(t, observable.Any(observable.False(), observable.False(), observable.True()))
	testspy.ExpectFail(t, observable.Any(observable.False(), observable.False()))

	testspy.ExpectPass(t, observable.All(observable.True(), observable.True(), observable.True()))
	testspy.ExpectFail(t, observable.All(observable.False(), observable.True(), observable.True()))

	testspy.ExpectFail(t, observable.Not(observable.Any)(observable.True(), observable.False()))
	testspy.ExpectFail(t, observable.Not(observable.All)(observable.True(), observable.True()))

	testspy.ExpectPass(t, observable.Not(observable.All)(observable.False(), observable.True(), observable.True()))
	testspy.ExpectPass(t, observable.Not(observable.Any)(observable.False(), observable.False(), observable.False()))
}

func TestNotVariadic(t *testing.T) {
	f := func(_ string, _ ...any) observable.Predicate {
		return observable.False()
	}
	testspy.ExpectPass(t, observable.Not(f)("foo", 1, 2, 3))
	testspy.ExpectFail(t, observable.Not(observable.Not(f))("foo", 1, 2, 3))
}
