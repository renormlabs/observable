// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable_test

import (
	"testing"

	"renorm.dev/observable"
	"renorm.dev/observable/internal/testspy"
)

func TestMapAsserts(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	testspy.ExpectPass(t, observable.ContainsKey(m, "a"))
	testspy.ExpectFail(t, observable.ContainsKey(m, "c"))
	testspy.ExpectPass(t, observable.ContainsValue(m, 1))
	testspy.ExpectFail(t, observable.ContainsValue(m, 3))

	testspy.ExpectPass(t, observable.MapLength(m, 2))
	testspy.ExpectFail(t, observable.MapLength(m, 1))

	newmap := map[string]int{
		"a": 1,
		"b": 2,
	}
	testspy.ExpectPass(t, observable.MapEqual(m, newmap))
	newmap["c"] = 3
	testspy.ExpectFail(t, observable.MapEqual(m, newmap))
}
