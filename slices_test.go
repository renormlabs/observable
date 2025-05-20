// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable_test

import (
	"testing"

	"renorm.dev/observable"
	"renorm.dev/observable/internal/testspy"
)

func TestSliceAsserts(t *testing.T) {
	foo := []int{1, 2, 3, 4, 5}

	testspy.ExpectPass(t, observable.Contains(foo, 2))
	testspy.ExpectFail(t, observable.Contains(foo, 7))
	testspy.ExpectPass(t, observable.Length(foo, 5))
	testspy.ExpectFail(t, observable.Length(foo, 2))

	testspy.ExpectPass(t, observable.SequenceEqual(foo, []int{1, 2, 3, 4, 5}))
	testspy.ExpectFail(t, observable.SequenceEqual(foo, []int{1, 2, 3, 4, 6}))
	testspy.ExpectFail(t, observable.SequenceEqual(foo, []int{1, 2, 3, 4}))
	testspy.ExpectPass(t, observable.ElementsMatch(foo, []int{2, 1, 3, 4, 5}))
	testspy.ExpectFail(t, observable.ElementsMatch(foo, []int{1, 2}))
	testspy.ExpectPass(t, observable.SequenceDeepEqual(foo, []int{1, 2, 3, 4, 5}))
	testspy.ExpectFail(t, observable.SequenceDeepEqual(foo, []int{1, 2, 3, 9, 5}))

	testspy.ExpectFail(t, observable.Empty(foo))
}
