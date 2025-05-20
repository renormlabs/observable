// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable_test

import (
	"testing"

	"renorm.dev/observable"
	"renorm.dev/observable/internal/testspy"
)

func TestChannelAsserts(t *testing.T) {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	testspy.ExpectPass(t, observable.ChanLength(ch, 2))
	testspy.ExpectFail(t, observable.ChanLength(ch, 5))
}
