package observable_test

import (
	"testing"

	"renorm.dev/observable"
	"renorm.dev/observable/internal/testspy"
)

func TestErrorIsChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.ErrorIs(errFoo, errFoo))
	testspy.ExpectFail(t, observable.ErrorIs(errFoo, errBar))

	testspy.ExpectPass(t, observable.Not(observable.ErrorIs)(errFoo, errBar))
	testspy.ExpectFail(t, observable.Not(observable.ErrorIs)(errFoo, errFoo))
}

func TestErrorsChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.Errors(func() error { return errFoo }))
	testspy.ExpectFail(t, observable.Errors(func() error { return nil }))

	testspy.ExpectPass(t, observable.Not(observable.Errors)(func() error { return nil }))
	testspy.ExpectFail(t, observable.Not(observable.Errors)(func() error { return errFoo }))
}

func TestErrorsWithChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.ErrorsWith(func() error { return errFoo }, errFoo))
	testspy.ExpectFail(t, observable.ErrorsWith(func() error { return errFoo }, errBar))

	testspy.ExpectPass(t, observable.Not(observable.ErrorsWith)(func() error { return errFoo }, errBar))
	testspy.ExpectFail(t, observable.Not(observable.ErrorsWith)(func() error { return errFoo }, errFoo))
}

func TestPanicsChecks(t *testing.T) {
	testspy.ExpectPass(t, observable.Panics(func() { panic("boom") }))
	testspy.ExpectFail(t, observable.Panics(func() {}))

	testspy.ExpectPass(t, observable.Not(observable.Panics)(func() {}))
	testspy.ExpectFail(t, observable.Not(observable.Panics)(func() { panic("boom") }))
}
