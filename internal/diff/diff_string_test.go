package diff_test

import (
	"strings"
	"testing"

	"github.com/seiyab/teq/internal/diff"
)

func TestDiff_MultilineString(t *testing.T) {
	t.Run("multi line string", func(t *testing.T) {
		type testCase struct {
			name  string
			left  string
			right string
			want  string
		}
		for _, tc := range []testCase{
			{
				name: "completely different",
				left: strings.Join([]string{
					"abc",
					"def",
				}, "\n"),
				right: strings.Join([]string{
					"jkl",
					"mno",
				}, "\n"),
				want: strings.Join([]string{
					`  string(`,
					`-   abc`,
					`-   def`,
					`+   jkl`,
					`+   mno`,
					`  )`,
				}, "\n"),
			},
			{
				name: "partial match",
				left: strings.Join([]string{
					"abc",
					"def",
					"ghi",
					"jkl",
				}, "\n"),
				right: strings.Join([]string{
					"abc",
					"ghi",
					"jkl",
					"mno",
				}, "\n"),
				want: strings.Join([]string{
					`  string(`,
					`    abc`,
					`-   def`,
					`    ghi`,
					`    jkl`,
					`+   mno`,
					`  )`,
				}, "\n"),
			},
			{
				name: "left is empty",
				left: "",
				right: strings.Join([]string{
					"abc",
					"ghi",
				}, "\n"),
				want: strings.Join([]string{
					`- ""`,
					`+ "abc\nghi"`,
				}, "\n"),
			},
			{
				name: "long a bit",
				left: strings.Join([]string{
					"---",
					"abc",
					"def",
					"ghi",
					"jkl",
					"mno",
					"pqr",
					"stu",
				}, "\n"),
				right: strings.Join([]string{
					"+++",
					"abc",
					"def",
					"ghi",
					"abc",
					"def",
					"pqr",
					"stu",
					"vwx",
				}, "\n"),
				want: strings.Join([]string{
					`  string(`,
					`-   ---`,
					`+   +++`,
					`    abc`,
					`    def`,
					`    ghi`,
					`-   jkl`,
					`-   mno`,
					`+   abc`,
					`+   def`,
					`    pqr`,
					`    stu`,
					`+   vwx`,
					`  )`,
				}, "\n"),
			},
			{
				name: "parial",
				left: strings.Join([]string{
					"---",
					"abc",
					"def",
				}, "\n"),
				right: strings.Join([]string{
					"+++",
					"def",
					"ghi",
				}, "\n"),
				want: strings.Join([]string{
					`  string(`,
					`-   ---`,
					`-   abc`,
					`+   +++`,
					`    def`,
					`+   ghi`,
					`  )`,
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
					t.Log(tc.want)
					t.Log(f)
				}
			})
		}
	})
}
