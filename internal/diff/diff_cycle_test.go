package diff_test

import (
	"strings"
	"testing"
)

func TestDiff_Cycle(t *testing.T) {
	t.Run("cyclic struct", func(t *testing.T) {
		type Cyclic struct {
			Ref *Cyclic
		}
		v := Cyclic{}
		v.Ref = &v
		runTest(t, v, Cyclic{Ref: &Cyclic{Ref: &v}}, "")
	})

	t.Run("different cyclic structs", func(t *testing.T) {
		type Cyclic struct {
			Value int
			Ref   *Cyclic
		}

		v1 := Cyclic{}
		tail := &v1
		for i := 1; i < 3; i++ {
			next := Cyclic{Value: i}
			tail.Ref = &next
			tail = &next
		}
		tail.Ref = &v1

		v2 := Cyclic{}
		tail = &v2
		for i := 1; i < 4; i++ {
			next := Cyclic{Value: i}
			tail.Ref = &next
			tail = &next
		}
		tail.Ref = &v2

		runTest(t, v1, v2, strings.Join([]string{
			`  diff_test.Cyclic{`,
			`:`,
			`    Ref: &diff_test.Cyclic{`,
			`:`,
			`      Ref: &diff_test.Cyclic{`,
			`        Value: 2,`,
			`        Ref: &diff_test.Cyclic{`,
			`-         Value: 0,`,
			`+         Value: 3,`,
			`          Ref: &diff_test.Cyclic{`,
			`-           Value: 1,`,
			`+           Value: 0,`,
			`-           Ref: &<circular reference>,`,
			`+           Ref: &<circular reference>,`,
			`          },`,
			`        },`,
			`      },`,
			`    },`,
			`  }`,
		}, "\n"))
	})
}
