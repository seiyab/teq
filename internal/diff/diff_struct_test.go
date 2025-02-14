package diff_test

import (
	"strings"
	"testing"
)

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
				runTest(t, tc.left, tc.right, tc.want)
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
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})
}
