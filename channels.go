// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable

import "fmt"

// ChanLength returns a [Predicate] that is ok when len(c) == want (buffered channels only).
func ChanLength[T any](c chan T, want int) Predicate {
	return Predicate{
		ok:  func() bool { return len(c) == want },
		msg: func() string { return fmt.Sprintf("expected channel buffer length %d, got %d", want, len(c)) },
	}
}
