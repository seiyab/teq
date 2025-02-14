package diff_test

import (
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
