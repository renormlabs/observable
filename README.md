# observable
Lightweight, zero-dependency helpers for writing expressive assertions in Go tests.

## Example

```go
package example_test

import (
    "errors"
    "testing"

    . "github.com/renormlabs/observable"
)

func TestAdd(t *testing.T) {
    add := func(a, b int) int { return a + b }

    // Basic equality
    Assert(t, Equal(add(2, 3), 5))

    // Negation
    Assert(t, Not(Equal)(add(2, 2), 5))

    // Error handling
    returnsErr := func() error { return errors.New("boom") }
    obs.Assert(t, Errors(returnsErr))
}
```
