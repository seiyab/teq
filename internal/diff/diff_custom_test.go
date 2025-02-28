package diff_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/seiyab/teq/internal/diff"
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

var jst *time.Location

func TestDiff_Stringer(t *testing.T) {
	t.Run("Stringer", func(t *testing.T) {
		runTest(t,
			time.Date(2025, 2, 3, 23, 3, 15, 0, time.UTC),
			time.Date(2024, 12, 19, 5, 45, 50, 0, jst),
			strings.Join([]string{
				`- time.Time("2025-02-03 23:03:15 +0000 UTC")`,
				`+ time.Time("2024-12-19 05:45:50 +0900 JST")`,
			}, "\n"),
		)
	})

	t.Run("Stringer in slice", func(t *testing.T) {
		runTest(t,
			[]time.Time{
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
			},
			strings.Join([]string{
				`  []time.Time{`,
				`:`,
				`    time.Time("2025-01-03 00:00:00 +0000 UTC"),`,
				`    time.Time("2025-01-04 00:00:00 +0000 UTC"),`,
				`-   time.Time("2025-01-05 00:00:00 +0000 UTC"),`,
				`    time.Time("2025-01-06 00:00:00 +0000 UTC"),`,
				`    time.Time("2025-01-07 00:00:00 +0000 UTC"),`,
				`:`,
				`  }`,
			}, "\n"),
		)
	})
}

func TestDiff_Format(t *testing.T) {
	o := diff.WithFormat(func(v int) string {
		return fmt.Sprintf("custom format(%d)", v)
	})

	t.Run("format", func(t *testing.T) {
		runTest(t,
			1,
			2,
			strings.Join([]string{
				`- int("custom format(1)")`,
				`+ int("custom format(2)")`,
			}, "\n"),
			o,
		)
	})

	t.Run("slice", func(t *testing.T) {
		runTest(t,
			[]int{1, 2, 3, 4},
			[]int{1, 2, 4},
			strings.Join([]string{
				`  []int{`,
				`    int("custom format(1)"),`,
				`    int("custom format(2)"),`,
				`-   int("custom format(3)"),`,
				`    int("custom format(4)"),`,
				`  }`,
			}, "\n"),
			o,
		)
	})
}

func TestDiff_MapKey(t *testing.T) {
	t.Run("Stringer", func(t *testing.T) {
		runTest(t,
			map[time.Duration]int{
				time.Second: 10,
				time.Minute: 2,
				time.Hour:   1,
			},
			map[time.Duration]int{
				time.Millisecond: 10,
				time.Minute:      3,
				time.Hour:        1,
			},
			strings.Join([]string{
				`  map[time.Duration]int{`,
				`+   time.Duration("1ms"): 10,`,
				`-   time.Duration("1s"): 10,`,
				`-   time.Duration("1m0s"): 2,`,
				`+   time.Duration("1m0s"): 3,`,
				`    time.Duration("1h0m0s"): 1,`,
				`  }`,
			}, "\n"),
		)
	})

	t.Run("Format", func(t *testing.T) {
		o := diff.WithFormat(func(v int) string {
			return fmt.Sprintf("custom format(%d)", v)
		})
		runTest(t,
			map[int]int{
				1: 10,
				2: 2,
				3: 1,
			},
			map[int]int{
				1: 10,
				2: 3,
				3: 1,
			},
			strings.Join([]string{
				`  map[int]int{`,
				`    int("custom format(1)"): int("custom format(10)"),`,
				`-   int("custom format(2)"): int("custom format(2)"),`,
				`+   int("custom format(2)"): int("custom format(3)"),`,
				`    int("custom format(3)"): int("custom format(1)"),`,
				`  }`,
			}, "\n"),
			o,
		)
	})
}
