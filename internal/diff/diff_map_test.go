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
					`    a: 1,`,
					`-   b: 2,`,
					`+   b: 3,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "different keys",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "c": 3},
				want: strings.Join([]string{
					`  map[string]int{`,
					`    a: 1,`,
					`-   b: 2,`,
					`+   c: 3,`,
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
					`+   a: 1,`,
					`+ }`,
				}, "\n"),
			},
			{
				name:  "nil maps (right)",
				left:  map[string]int{"a": 1},
				right: nil,
				want: strings.Join([]string{
					`- map[string]int{`,
					`-   a: 1,`,
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
}
