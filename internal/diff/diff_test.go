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

func TestDiff_Struct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type S struct {
			A int
			B string
		}
		type testCase struct {
			name  string
			left  S
			right S
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "completely different",
				left:  S{A: 1, B: "hello"},
				right: S{A: 2, B: "world"},
				want: strings.Join([]string{
					`  diff_test.S{`,
					`-   A: 1`,
					`+   A: 2`,
					`-   B: "hello"`,
					`+   B: "world"`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "partially different",
				left:  S{A: 1, B: "hello"},
				right: S{A: 1, B: "world"},
				want: strings.Join([]string{
					`  diff_test.S{`,
					`    A: 1`,
					`-   B: "hello"`,
					`+   B: "world"`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				d, e := diff.New().Diff(tc.left, tc.right)
				if e != nil {
					t.Fatal(e)
				}
				f := d.Format()
				if f != tc.want {
					t.Errorf("expected %q, got %q", tc.want, f)
					t.Log(f)
					t.Log(tc.want)
				}
			})
		}
	})
}
