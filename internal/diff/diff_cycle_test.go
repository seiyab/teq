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

	t.Run("cyclic slice", func(t *testing.T) {
		t.Run("case 1", func(t *testing.T) {
			v1 := []interface{}{1, 2, nil}
			v1[2] = v1
			v2 := []interface{}{1, 2, nil, nil}
			v2[2], v2[3] = v2, v2

			runTest(t, v1, v2, strings.Join([]string{
				"  []interface {}{",
				"    1,",
				"    2,",
				"+   []interface {}{",
				"+     1,",
				"+     2,",
				"+     <circular reference>,",
				"+     []interface {}{",
				"+       1,",
				"+       2,",
				"+       <circular reference>,",
				"+       <circular reference>,",
				"+     },",
				"+   },",
				"    []interface {}{",
				"      1,",
				":",
				"      []interface {}{",
				"        1,",
				"        2,",
				"-       <circular reference>,",
				"+       <circular reference>,",
				"+       <circular reference>,",
				"      },",
				"+     <circular reference>,",
				"    },",
				"  }",
			}, "\n"))
		})

		t.Run("case 2", func(t *testing.T) {
			v1 := []interface{}{1, nil}
			v1[1] = v1
			v2 := []interface{}{2, nil}
			v2[1] = v2

			runTest(t, v1, v2, strings.Join([]string{
				`  []interface {}{`,
				`-   1,`,
				`+   2,`,
				`    []interface {}{`,
				`-     1,`,
				`-     <circular reference>,`,
				`+     2,`,
				`+     <circular reference>,`,
				`    },`,
				`  }`,
			}, "\n"))
		})

		t.Run("case 3", func(t *testing.T) {
			v1 := []interface{}{nil, 1}
			v1[0] = v1
			v2 := []interface{}{nil, 2, 1}
			v2[0] = v2

			runTest(t, v1, v2, strings.Join([]string{
				`  []interface {}{`,
				`    []interface {}{`,
				`-     <circular reference>,`,
				`+     <circular reference>,`,
				`+     2,`,
				`      1,`,
				`    },`,
				`+   2,`,
				`    1,`,
				`  }`,
			}, "\n"))
		})
	})

	t.Run("cyclic interface", func(t *testing.T) {
		type Cyclic struct {
			Value int
			Ref   interface{}
		}
		var v1, v2 Cyclic
		v1 = Cyclic{
			Value: 1,
			Ref:   &v1,
		}
		v2 = Cyclic{
			Value: 1,
			Ref: &Cyclic{
				Value: 2,
				Ref:   &v2,
			},
		}
		runTest(t, v1, v2, strings.Join([]string{
			`  diff_test.Cyclic{`,
			`    Value: 1,`,
			`    Ref: &diff_test.Cyclic{`,
			`-     Value: 1,`,
			`+     Value: 2,`,
			`      Ref: &diff_test.Cyclic{`,
			`        Value: 1,`,
			`-       Ref: &<circular reference>,`,
			`+       Ref: &<circular reference>,`,
			`      },`,
			`    },`,
			`  }`,
		}, "\n"),
		)
	})
}
