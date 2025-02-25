package diff_test

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/seiyab/teq/internal/diff"
)

func TestDiff_String(t *testing.T) {
	t.Run("word", func(t *testing.T) {
		runTest(t, "hello", "world", strings.Join([]string{
			`- "hello"`,
			`+ "world"`,
		}, "\n"))
	})
}

func TestDiff_Primitive(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		runTest(t, 1, 2, strings.Join([]string{
			`- 1`,
			`+ 2`,
		}, "\n"))
	})

	t.Run("float", func(t *testing.T) {
		runTest(t, 1.0, 2.0, strings.Join([]string{
			`- 1.000000`,
			`+ 2.000000`,
		}, "\n"))
	})

	t.Run("bool", func(t *testing.T) {
		runTest(t, true, false, strings.Join([]string{
			`- true`,
			`+ false`,
		}, "\n"))
	})
}

func TestDiff_Func(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		// NOTE: reflect.DeepEqual cannot compare functions.
		runTest(t, func(*testing.T) {}, func(*testing.T) {}, strings.Join([]string{
			`- func(*testing.T) { ... }`,
			`+ func(*testing.T) { ... }`,
		}, "\n"))
	})

	type f = func()
	t.Run("both nil", func(t *testing.T) {
		runTest(t, f(nil), f(nil), "")
	})

	t.Run("left nil", func(t *testing.T) {
		t.Skip()
		expected := strings.Join([]string{
			`- func(t *testing.T) { ... }`,
			`+ func()(nil)`,
		}, "\n")
		runTest(t, TestDiff_Primitive, f(nil), expected)
	})
}

func TestDiff_Pointer(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		runTest(t, (*int)(nil), (*int)(nil), "")
	})

	t.Run("same", func(t *testing.T) {
		runTest(t, ref(100), ref(100), "")
	})

	t.Run("different", func(t *testing.T) {
		expected := strings.Join([]string{
			`- &100`,
			`+ &200`,
		}, "\n")
		runTest(t, ref(100), ref(200), expected)
	})
}

func ref[T any](x T) *T {
	return &x
}

func TestDiff_Interface(t *testing.T) {
	type S struct {
		X fmt.Stringer
		Y io.Closer
	}

	t.Run("same", func(t *testing.T) {
	})

	t.Run("different Stringer", func(t *testing.T) {
		runTest(t,
			S{X: I(1)},
			S{X: I(2)},
			strings.Join([]string{
				`  diff_test.S{`,
				`-   X: diff_test.I("1"),`,
				`+   X: diff_test.I("2"),`,
				`    Y: io.Closer(nil),`,
				`  }`,
			}, "\n"))
	})

	t.Run("different", func(t *testing.T) {
		runTest(t,
			S{X: I(1), Y: C(1)},
			S{X: I(1), Y: C(2)},
			strings.Join([]string{
				`  diff_test.S{`,
				`    X: diff_test.I("1"),`,
				`-   Y: 1,`,
				`+   Y: 2,`,
				`  }`,
			}, "\n"))
	})

	t.Run("nil", func(t *testing.T) {
		runTest(t,
			S{X: I(1)},
			S{X: nil},
			strings.Join([]string{
				`  diff_test.S{`,
				`-   X: diff_test.I("1"),`,
				`+   X: fmt.Stringer(nil),`,
				`    Y: io.Closer(nil),`,
				`  }`,
			}, "\n"))
	})
}

type I int

func (v I) String() string { return strconv.Itoa(int(v)) }

type C int

func (v C) Close() error { return nil }

func TestDiff_Chan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		type s struct {
			c chan bool
		}
		runTest(t, s{c: nil}, s{c: make(chan bool)}, strings.Join([]string{
			`  diff_test.s{`,
			`-   c: chan bool(nil),`,
			`+   c: chan bool,`,
			`  }`,
		}, "\n"))
	})

	t.Run("different", func(t *testing.T) {
		expected := strings.Join([]string{
			`- <-chan string`,
			`+ <-chan string`,
		}, "\n")
		runTest(t, make(<-chan string), make(<-chan string), expected)
	})
}

func TestDiff_TypeMismatch(t *testing.T) {
	t.Run("struct vs struct", func(t *testing.T) {
		type s struct {
			i int
		}
		type u struct {
			i int
		}
		runTest(t, s{1}, u{1}, strings.Join([]string{
			`- diff_test.s{`,
			`-   i: 1,`,
			`- }`,
			`+ diff_test.u{`,
			`+   i: 1,`,
			`+ }`,
		}, "\n"))
	})

	t.Run("string vs map", func(t *testing.T) {
		runTest(t, "hello", map[string]int{}, strings.Join([]string{
			`- "hello"`,
			`+ map[string]int{`,
			`+ }`,
		}, "\n"))
	})
}

func runTest(t *testing.T, left, right any, want string, opts ...diff.Option) {
	t.Helper()
	d := diff.DiffString(left, right, opts...)
	if d != want {
		t.Errorf("expected %q, got %q", want, d)
		p := diff.DiffString(want, d)
		t.Log("\n" + p)
	}
}
