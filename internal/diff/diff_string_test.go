package diff_test

import (
	"strings"
	"testing"
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
			{
				name: "sight width",
				left: strings.Join([]string{
					"abc",
					"def",
					"ghi",
					"jkl",
					"mno",
					"pqr",
					"stu",
					"vwx",
					"yz",
				}, "\n"),
				right: strings.Join([]string{
					"abc",
					"def",
					"ghi",
					"jkl",
					"pqr",
					"stu",
					"vwx",
					"yz",
				}, "\n"),
				want: strings.Join([]string{
					`  string(`,
					`:`,
					`    ghi`,
					`    jkl`,
					`-   mno`,
					`    pqr`,
					`    stu`,
					`:`,
					`  )`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})
}
