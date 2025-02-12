package diff_test

import (
	"strings"
	"testing"

	"github.com/seiyab/teq/internal/diff"
)

func TestDiff_Map(t *testing.T) {
	t.Run("map of primitive", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[string]int
			right map[string]int
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "identical",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "b": 2},
				want:  "",
			},
			{
				name:  "different values",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "b": 3},
				want: strings.Join([]string{
					`  map[string]int{`,
					`    "a": 1,`,
					`-   "b": 2,`,
					`+   "b": 3,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "different keys",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "c": 3},
				want: strings.Join([]string{
					`  map[string]int{`,
					`    "a": 1,`,
					`-   "b": 2,`,
					`+   "c": 3,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "nil maps",
				left:  nil,
				right: map[string]int{"a": 1},
				want: strings.Join([]string{
					`- map[string]int(nil)`,
					`+ map[string]int{`,
					`+   "a": 1,`,
					`+ }`,
				}, "\n"),
			},
			{
				name:  "nil maps (right)",
				left:  map[string]int{"a": 1},
				right: nil,
				want: strings.Join([]string{
					`- map[string]int{`,
					`-   "a": 1,`,
					`- }`,
					`+ map[string]int(nil)`,
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
					p, e := diff.New().Diff(tc.want, f)
					if e != nil {
						t.Fatal(e)
					}
					t.Log(p.Format())
				}
			})
		}
	})

	t.Run("map with int key", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[int]string
			right map[int]string
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "identical",
				left:  map[int]string{1: "a", 2: "b"},
				right: map[int]string{1: "a", 2: "b"},
				want:  "",
			},
			{
				name:  "different values",
				left:  map[int]string{1: "a", 2: "b"},
				right: map[int]string{1: "a", 2: "c"},
				want: strings.Join([]string{
					`  map[int]string{`,
					`    1: "a",`,
					`-   2: "b",`,
					`+   2: "c",`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runMapTest(t, tc.left, tc.right, tc.want)
			})
		}
	})

	t.Run("nested map", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[string]map[string]int
			right map[string]map[string]int
			want  string
		}
		for _, tc := range []testCase{
			{
				name: "different nested values",
				left: map[string]map[string]int{
					"x": {"a": 1, "b": 2},
					"y": {"c": 3},
				},
				right: map[string]map[string]int{
					"x": {"a": 1, "b": 3},
					"y": {"c": 3},
				},
				want: strings.Join([]string{
					`  map[string]map[string]int{`,
					`    "x": map[string]int{`,
					`      "a": 1,`,
					`-     "b": 2,`,
					`+     "b": 3,`,
					`    },`,
					`    "y": map[string]int{`,
					`:`,
					`  }`,
				}, "\n"),
			},
			{
				name: "nil nested map",
				left: map[string]map[string]int{
					"x": {"a": 1},
					"y": nil,
				},
				right: map[string]map[string]int{
					"x": {"a": 1},
					"y": {"b": 2},
				},
				want: strings.Join([]string{
					`  map[string]map[string]int{`,
					`:`,
					`      "a": 1,`,
					`    },`,
					`-   "y": map[string]int(nil),`,
					`+   "y": map[string]int{`,
					`+     "b": 2,`,
					`+   },`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runMapTest(t, tc.left, tc.right, tc.want)
			})
		}
	})

	t.Run("map with interface key", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[any]int
			right map[any]int
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "mixed key types",
				left:  map[any]int{1: 10, "x": 20, true: 30},
				right: map[any]int{1: 10, "x": 25, true: 30},
				want: strings.Join([]string{
					`  map[interface {}]int{`,
					`-   "x": 20,`,
					`+   "x": 25,`,
					`    1: 10,`,
					`    true: 30,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "with nil key",
				left:  map[any]int{nil: 1, "x": 2},
				right: map[any]int{nil: 1, "y": 2},
				want: strings.Join([]string{
					`  map[interface {}]int{`,
					`-   "x": 2,`,
					`+   "y": 2,`,
					`    interface {}(<nil>): 1,`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runMapTest(t, tc.left, tc.right, tc.want)
			})
		}
	})
}

func runMapTest(t *testing.T, left, right any, want string) {
	t.Helper()
	d, e := diff.New().Diff(left, right)
	if e != nil {
		t.Fatal(e)
	}
	f := d.Format()
	if f != want {
		t.Errorf("expected %q, got %q", want, f)
		p, e := diff.New().Diff(want, f)
		if e != nil {
			t.Fatal(e)
		}
		t.Log(p.Format())
	}
}
