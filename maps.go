package observable

import (
	"fmt"
	"reflect"
	"sync"
)

// ContainsKey succeeds when key exists in map m.
func ContainsKey[K comparable, V any](m map[K]V, key K) Predicate {
	return Predicate{
		ok:  func() bool { _, ok := m[key]; return ok },
		msg: func() string { return fmt.Sprintf("expected map to contain key %v", key) },
	}
}

// ContainsValue succeeds when val appears as a value in map m.
func ContainsValue[K comparable, V comparable](m map[K]V, val V) Predicate {
	return Predicate{
		ok: func() bool {
			for _, v := range m {
				if v == val {
					return true
				}
			}
			return false
		},
		msg: func() string { return fmt.Sprintf("expected map to contain value %v", val) },
	}
}

// MapEqual succeeds when two maps are deeply equal (reflect.DeepEqual).
func MapEqual[K comparable, V any](got, want map[K]V) Predicate {
	var (
		once  sync.Once
		match bool
	)

	check := func() { once.Do(func() { match = reflect.DeepEqual(got, want) }) }

	return Predicate{
		ok: func() bool {
			check()
			return match
		},
		msg: func() string {
			check()
			return fmt.Sprintf("expected maps to be equal\nwant: %#v\ngot:  %#v", want, got)
		},
	}
}

// MapLength succeeds when len(m) == want.
func MapLength[K comparable, V any](m map[K]V, want int) Predicate {
	return Predicate{
		ok:  func() bool { return len(m) == want },
		msg: func() string { return fmt.Sprintf("expected map size %d, got %d", want, len(m)) },
	}
}
