package diff_test

import (
	"fmt"
	"strings"
	"testing"
)

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
				want:  "",
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
				runTest(t, tc.left, tc.right, tc.want)
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
				want:  "",
			},
			{
				left:  [][]int{{1, 2}, {3, 4}},
				right: [][]int{{1, 2}, {3, 5}},
				want: strings.Join([]string{
					`  [][]int{`,
					`:`,
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
					`:`,
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
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})
}

func TestDiff_Array(t *testing.T) {
	t.Run("array of primitive", func(t *testing.T) {
		type testCase struct {
			left  [3]int
			right [3]int
			want  string
		}
		for _, tc := range []testCase{
			{
				left:  [3]int{1, 2, 3},
				right: [3]int{1, 2, 3},
				want:  "",
			},
			{
				left:  [3]int{1, 2, 3},
				right: [3]int{1, 2, 4},
				want: strings.Join([]string{
					`  [3]int{`,
					`    1,`,
					`    2,`,
					`-   3,`,
					`+   4,`,
					`  }`,
				}, "\n"),
			},
			{
				left:  [3]int{1, 2, 3},
				right: [3]int{2, 3, 4},
				want: strings.Join([]string{
					`  [3]int{`,
					`-   1,`,
					`    2,`,
					`    3,`,
					`+   4,`,
					`  }`,
				}, "\n"),
			},
		} {
			name := fmt.Sprintf("%v vs %v", tc.left, tc.right)
			t.Run(name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})

	t.Run("array of array", func(t *testing.T) {
		type testCase struct {
			left  [2][2]int
			right [2][2]int
			want  string
		}
		for _, tc := range []testCase{
			{
				left:  [2][2]int{{1, 2}, {3, 4}},
				right: [2][2]int{{1, 2}, {3, 4}},
				want:  "",
			},
			{
				left:  [2][2]int{{1, 2}, {3, 4}},
				right: [2][2]int{{1, 2}, {3, 5}},
				want: strings.Join([]string{
					`  [2][2]int{`,
					`:`,
					`    [2]int{`,
					`      3,`,
					`-     4,`,
					`+     5,`,
					`    },`,
					`  }`,
				}, "\n"),
			},
			{
				left:  [2][2]int{{1, 2}, {3, 4}},
				right: [2][2]int{{1, 3}, {2, 4}},
				want: strings.Join([]string{
					`  [2][2]int{`,
					`    [2]int{`,
					`      1,`,
					`-     2,`,
					`+     3,`,
					`    },`,
					`    [2]int{`,
					`-     3,`,
					`+     2,`,
					`      4,`,
					`    },`,
					`  }`,
				}, "\n"),
			},
		} {
			name := fmt.Sprintf("%v vs %v", tc.left, tc.right)
			t.Run(name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})
}
