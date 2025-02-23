package diff_test

import (
	"strings"
	"testing"
	"time"
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

var jst *time.Location

func TestDiff_MarshalText(t *testing.T) {
	t.Run("TextMarshaler", func(t *testing.T) {
		runTest(t,
			time.Date(2025, 2, 3, 23, 3, 15, 0, time.UTC),
			time.Date(2024, 12, 19, 5, 45, 50, 0, jst),
			strings.Join([]string{
				`- time.Time("2025-02-03T23:03:15Z")`,
				`+ time.Time("2024-12-19T05:45:50+09:00")`,
			}, "\n"),
		)
	})

	t.Run("TextMarshaler in slice", func(t *testing.T) {
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
				`    time.Time("2025-01-03T00:00:00Z"),`,
				`    time.Time("2025-01-04T00:00:00Z"),`,
				`-   time.Time("2025-01-05T00:00:00Z"),`,
				`    time.Time("2025-01-06T00:00:00Z"),`,
				`    time.Time("2025-01-07T00:00:00Z"),`,
				`:`,
				`  }`,
			}, "\n"),
		)
	})
}
