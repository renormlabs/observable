// Copyright (c) 2025 Renorm Labs. All rights reserved.

package observable

import (
	"errors"
	"fmt"
)

// ErrorIs returns a [Predicate] that is ok when [errors.Is](err, target) is true.
func ErrorIs(err, target error) Predicate {
	return Predicate{
		ok: func() bool { return errors.Is(err, target) },
		msg: func() string {
			return fmt.Sprintf("expected error %v to match %v", err, target)
		},
	}
}

// Errors returns a [Predicate] that is ok when f returns a non‑nil error.
func Errors(f func() error) Predicate {
	return Predicate{
		ok:  func() bool { return f() != nil },
		msg: func() string { return "expected function to return a non-nil error" },
	}
}

// ErrorsWith returns a [Predicate] that is ok when f returns an error that matches target according to [errors.Is].
func ErrorsWith(f func() error, target error) Predicate {
	return Predicate{
		ok:  func() bool { return errors.Is(f(), target) },
		msg: func() string { return fmt.Sprintf("expected returned error to match %v", target) },
	}
}

// Panics returns a [Predicate] that is ok when f panics.
func Panics(f func()) Predicate {
	return Predicate{
		ok: func() (panicked bool) {
			defer func() {
				if recover() != nil {
					panicked = true
				}
			}()
			f()
			return
		},
		msg: func() string { return "expected function to panic" },
	}
}
