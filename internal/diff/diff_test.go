package diff_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/seiyab/teq/internal/diff"
)

func TestDiff_String(t *testing.T) {
	t.Run("word", func(t *testing.T) {
		d, e := diff.New().Diff("hello", "world")
		if e != nil {
			t.Fatal(e)
		}
		f := d.Format()
		x := strings.Join([]string{
			`- "hello"`,
			`+ "world"`,
		}, "\n")
		if f != x {
			t.Fatalf("expected %q, got %q", x, f)
		}
	})
}

func TestDiff_Primitive(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		d, e := diff.New().Diff(1, 2)
		if e != nil {
			t.Fatal(e)
		}
		f := d.Format()
		x := strings.Join([]string{
			`- 1`,
			`+ 2`,
		}, "\n")
		if f != x {
			t.Fatalf("expected %q, got %q", x, f)
		}
	})

	t.Run("float", func(t *testing.T) {
		d, e := diff.New().Diff(1.0, 2.0)
		if e != nil {
			t.Fatal(e)
		}
		f := d.Format()
		x := strings.Join([]string{
			`- 1.000000`,
			`+ 2.000000`,
		}, "\n")
		if f != x {
			t.Fatalf("expected %q, got %q", x, f)
		}
	})

	t.Run("bool", func(t *testing.T) {
		d, e := diff.New().Diff(true, false)
		if e != nil {
			t.Fatal(e)
		}
		f := d.Format()
		x := strings.Join([]string{
			`- true`,
			`+ false`,
		}, "\n")
		if f != x {
			t.Fatalf("expected %q, got %q", x, f)
		}
	})
}

func TestDiff_Func(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		// NOTE: reflect.DeepEqual cannot compare functions.
		expected := strings.Join([]string{
			`- func(*testing.T) { ... }`,
			`+ func(*testing.T) { ... }`,
		}, "\n")
		runTest(t, TestDiff_Func, TestDiff_Func, expected)
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
		E error
	}

	t.Run("same", func(t *testing.T) {
	})

	t.Run("different", func(t *testing.T) {
		runTest(t,
			S{X: I(1), E: fmt.Errorf("world")},
			S{X: I(2), E: fmt.Errorf("world")},
			strings.Join([]string{
				`  diff_test.S{`,
				`-   X: 1,`,
				`+   X: 2,`,
				`    E: &errors.errorString{`,
				`      s: "world",`,
				`:`,
				`  }`,
			}, "\n"))
	})

	t.Run("nil", func(t *testing.T) {
		runTest(t,
			S{X: I(1), E: fmt.Errorf("world")},
			S{X: nil, E: fmt.Errorf("world")},
			strings.Join([]string{
				`  diff_test.S{`,
				`-   X: 1,`,
				`+   X: fmt.Stringer(nil),`,
				`    E: &errors.errorString{`,
				`      s: "world",`,
				`:`,
				`  }`,
			}, "\n"))
	})
}

type I int

func (v I) String() string { return strconv.Itoa(int(v)) }

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

func runTest(t *testing.T, left, right any, want string) {
	t.Helper()
	d, err := diff.DiffString(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if d != want {
		t.Errorf("expected %q, got %q", want, d)
		p, e := diff.New().Diff(want, d)
		if e != nil {
			t.Fatal(e)
		}
		for _, l := range strings.Split(p.Format(), "\n") {
			t.Log(l)
		}
	}
}
