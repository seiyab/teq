package diff_test

import (
	"fmt"
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
					`-   A: 1,`,
					`+   A: 2,`,
					`-   B: "hello",`,
					`+   B: "world",`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "partially different",
				left:  S{A: 1, B: "hello"},
				right: S{A: 1, B: "world"},
				want: strings.Join([]string{
					`  diff_test.S{`,
					`    A: 1,`,
					`-   B: "hello",`,
					`+   B: "world",`,
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

	t.Run("nested struct", func(t *testing.T) {
		type S1 struct {
			A int
			B string
		}
		type S2 struct {
			X S1
			Y S1
		}
		type testCase struct {
			name  string
			left  S2
			right S2
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "completely different",
				left:  S2{X: S1{A: 1, B: "hello"}, Y: S1{A: 2, B: "world"}},
				right: S2{X: S1{A: 2, B: "world"}, Y: S1{A: 1, B: "hello"}},
				want: strings.Join([]string{
					`  diff_test.S2{`,
					`    X: diff_test.S1{`,
					`-     A: 1,`,
					`+     A: 2,`,
					`-     B: "hello",`,
					`+     B: "world",`,
					`    },`,
					`    Y: diff_test.S1{`,
					`-     A: 2,`,
					`+     A: 1,`,
					`-     B: "world",`,
					`+     B: "hello",`,
					`    },`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "partially different",
				left:  S2{X: S1{A: 1, B: "hello"}, Y: S1{A: 1, B: "world"}},
				right: S2{X: S1{A: 1, B: "world"}, Y: S1{A: 1, B: "hello"}},
				want: strings.Join([]string{
					`  diff_test.S2{`,
					`    X: diff_test.S1{`,
					`      A: 1,`,
					`-     B: "hello",`,
					`+     B: "world",`,
					`    },`,
					`    Y: diff_test.S1{`,
					`      A: 1,`,
					`-     B: "world",`,
					`+     B: "hello",`,
					`    },`,
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

func TestDiff_Slice(t *testing.T) {
	t.Run("slice of primitive", func(t *testing.T) {
		type testCase struct {
			left  []int
			right []int
			want  string
		}
		for _, tc := range []testCase{
			{
				left:  []int{1, 2, 3},
				right: []int{1, 2, 3},
				want: strings.Join([]string{
					`  []int{`,
					`    1,`,
					`    2,`,
					`    3,`,
					`  }`,
				}, "\n"),
			},
			{
				left:  []int{1, 2, 3},
				right: []int{1, 2, 4},
				want: strings.Join([]string{
					`  []int{`,
					`    1,`,
					`    2,`,
					`-   3,`,
					`+   4,`,
					`  }`,
				}, "\n"),
			},
			{
				left:  []int{1, 2, 3},
				right: []int{2, 3},
				want: strings.Join([]string{
					`  []int{`,
					`-   1,`,
					`    2,`,
					`    3,`,
					`  }`,
				}, "\n"),
			},
			{
				left:  []int{1, 5},
				right: []int{1, 2, 3, 4, 5},
				want: strings.Join([]string{
					`  []int{`,
					`    1,`,
					`+   2,`,
					`+   3,`,
					`+   4,`,
					`    5,`,
					`  }`,
				}, "\n"),
			},
			{
				left:  []int{1, 3},
				right: []int{},
				want: strings.Join([]string{
					`  []int{`,
					`-   1,`,
					`-   3,`,
					`  }`,
				}, "\n"),
			},
			{
				left:  []int{},
				right: []int{1, 2},
				want: strings.Join([]string{
					`  []int{`,
					`+   1,`,
					`+   2,`,
					`  }`,
				}, "\n"),
			},
		} {
			name := fmt.Sprintf("%v vs %v", tc.left, tc.right)
			t.Run(name, func(t *testing.T) {
				d, e := diff.New().Diff(tc.left, tc.right)
				if e != nil {
					t.Fatal(e)
				}
				f := d.Format()
				if f != tc.want {
					t.Errorf("expected %q, got %q", tc.want, f)
					t.Log(tc.want)
					t.Log(f)
				}
			})
		}
	})

	t.Run("slice of slice", func(t *testing.T) {
		type testCase struct {
			left  [][]int
			right [][]int
			want  string
		}
		for _, tc := range []testCase{
			{
				left:  [][]int{{1, 2}, {3, 4}},
				right: [][]int{{1, 2}, {3, 4}},
				want: strings.Join([]string{
					`  [][]int{`,
					`    []int{`,
					`      1,`,
					`      2,`,
					`    },`,
					`    []int{`,
					`      3,`,
					`      4,`,
					`    },`,
					`  }`,
				}, "\n"),
			},
			{
				left:  [][]int{{1, 2}, {3, 4}},
				right: [][]int{{1, 2}, {3, 5}},
				want: strings.Join([]string{
					`  [][]int{`,
					`    []int{`,
					`      1,`,
					`      2,`,
					`    },`,
					`    []int{`,
					`      3,`,
					`-     4,`,
					`+     5,`,
					`    },`,
					`  }`,
				}, "\n"),
			},
			{
				left:  [][]int{{1, 2}, {3, 4}},
				right: [][]int{{1, 2}},
				want: strings.Join([]string{
					`  [][]int{`,
					`    []int{`,
					`      1,`,
					`      2,`,
					`    },`,
					`-   []int{`,
					`-     3,`,
					`-     4,`,
					`-   },`,
					`  }`,
				}, "\n"),
			},
			{
				left:  [][]int{{1, 2}, {3, 4}},
				right: [][]int{{1, 3}, {2}, {4, 5}},
				want: strings.Join([]string{
					`  [][]int{`,
					`+   []int{`,
					`+     1,`,
					`+     3,`,
					`+   },`,
					`    []int{`,
					`-     1,`,
					`      2,`,
					`    },`,
					`    []int{`,
					`-     3,`,
					`      4,`,
					`+     5,`,
					`    },`,
					`  }`,
				}, "\n"),
			},
		} {
			name := fmt.Sprintf("%v vs %v", tc.left, tc.right)
			t.Run(name, func(t *testing.T) {
				d, e := diff.New().Diff(tc.left, tc.right)
				if e != nil {
					t.Fatal(e)
				}
				f := d.Format()
				if f != tc.want {
					t.Errorf("expected %q, got %q", tc.want, f)
					t.Log(tc.want)
					t.Log(f)
				}
			})
		}
	})
}
