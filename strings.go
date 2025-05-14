package observable

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

// StringLength returns a [Predicate] that succeeds when len(s) == want.
func StringLength(s string, want int) Predicate {
	return Predicate{
		ok:  func() bool { return len(s) == want },
		msg: func() string { return fmt.Sprintf("expected length %d, got %d", want, len(s)) },
	}
}

// EmptyString returns a [Predicate] that succeeds when the string is "".
func EmptyString(s string) Predicate {
	return StringLength(s, 0)
}

// RuneLength returns a [Predicate] that succeeds when utf8.RuneCountInString(s) == want.
func RuneLength(s string, want int) Predicate {
	return Predicate{
		ok:  func() bool { return utf8.RuneCountInString(s) == want },
		msg: func() string { return fmt.Sprintf("expected rune length %d, got %d", want, utf8.RuneCountInString(s)) },
	}
}

// HasPrefix returns a [Predicate] that succeeds when strings.HasPrefix(s, prefix).
func HasPrefix(s, prefix string) Predicate {
	return Predicate{
		ok:  func() bool { return strings.HasPrefix(s, prefix) },
		msg: func() string { return fmt.Sprintf("expected %q to have prefix %q", s, prefix) },
	}
}

// HasSuffix returns a [Predicate] that succeeds when strings.HasSuffix(s, suffix).
func HasSuffix(s, suffix string) Predicate {
	return Predicate{
		ok:  func() bool { return strings.HasSuffix(s, suffix) },
		msg: func() string { return fmt.Sprintf("expected %q to have suffix %q", s, suffix) },
	}
}

// ContainsSubstring returns a [Predicate] that succeeds when strings.Contains(s, substr).
func ContainsSubstring(s, substr string) Predicate {
	return Predicate{
		ok:  func() bool { return strings.Contains(s, substr) },
		msg: func() string { return fmt.Sprintf("expected %q to contain %q", s, substr) },
	}
}

// EqualFold returns a [Predicate] that succeeds when strings.EqualFold(got, want) (case-insensitive).
func EqualFold(got, want string) Predicate {
	return Predicate{
		ok:  func() bool { return strings.EqualFold(got, want) },
		msg: func() string { return fmt.Sprintf("expected %q (case-insensitive), got %q", want, got) },
	}
}

// RegexpMatches returns a [Predicate] that succeeds when the regular expression re matches s. The regular expression can either be a [*regexp.Regexp] or a string which will be compiled with [regexp.MustCompile].
func RegexpMatches[T reOrStringT](s string, reOrString T) Predicate {
	var (
		once sync.Once
		re   *regexp.Regexp
	)

	eval := func() {
		once.Do(func() {
			switch x := any(reOrString).(type) {
			case *regexp.Regexp:
				re = x
			case string:
				re = regexp.MustCompile(x)
			}
		})
	}

	return Predicate{
		ok:  func() bool { eval(); return re.MatchString(s) },
		msg: func() string { eval(); return fmt.Sprintf("expected %q to match %q", s, re.String()) },
	}
}

type reOrStringT interface {
	string | *regexp.Regexp
}
