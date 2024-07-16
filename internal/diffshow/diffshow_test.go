package diffshow_test

import (
	"strings"
	"testing"

	"github.com/seiyab/teq"
	"github.com/seiyab/teq/internal/diffshow"
)

func TestCompare(t *testing.T) {
	type testCase struct {
		name string
		a    any
		b    any
		diff string
	}
	testCases := []testCase{
		{
			name: "int",
			a:    1,
			b:    2,
			diff: join(
				`--- expected`,
				`+++ actual`,
				`@@ -1,2 +1,2 @@`,
				`-int(1)`,
				`+int(2)`,
				` `,
				``,
			),
		},
		{
			name: "inline string",
			a:    "Hello",
			b:    "World",
			diff: join(
				`--- expected`,
				`+++ actual`,
				`@@ -1,2 +1,2 @@`,
				`-"Hello"`,
				`+"World"`,
				` `,
				``,
			),
		},
	}
	tq := teq.New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := diffshow.Compare(tc.a, tc.b, tq.Eq)
			got, err := d.Format(tq.Fmt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := tc.diff
			tq.Equal(t, want, got)
			if got != want {
				t.Errorf("unexpected diff\n got: %s\nwant: %s", got, want)
			}
		})
	}
}

func join(a ...string) string {
	return strings.Join(a, "\n")
}
