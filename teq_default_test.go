package teq_test

import (
	"fmt"
	"testing"

	"github.com/seiyab/teq"
)

type test struct {
	a             any
	b             any
	expected      []string
	pendingFormat bool // for development. we don't have stable format yet.
}

type group struct {
	name  string
	tests []test
}

func TestEqual(t *testing.T) {
	assert := teq.New()

	groups := []group{
		{"primitives", primitives()},
		{"structs", structs()},
		{"slices", slices()},
		{"maps", maps()},
		{"recursions", recursions()},
	}

	for _, group := range groups {
		t.Run(group.name, func(t *testing.T) {
			for _, test := range group.tests {
				name := fmt.Sprintf("%T(%v) == %T(%v)", test.a, test.a, test.b, test.b)
				t.Run(name, func(t *testing.T) {
					mt := &mockT{}
					assert.Equal(mt, test.a, test.b)
					if len(mt.errors) != len(test.expected) {
						if len(mt.errors) > len(test.expected) {
							for _, e := range mt.errors {
								t.Logf("got %q", e)
							}
						}
						t.Fatalf("expected %d errors, got %d", len(test.expected), len(mt.errors))
					}

					if !test.pendingFormat {
						for i, e := range test.expected {
							if mt.errors[i] != e {
								t.Errorf("expected %q, got %q at i = %d", e, mt.errors[i], i)
							}
						}
					}

					{
						mt := &mockT{}
						assert.NotEqual(mt, test.a, test.b)
						if (len(mt.errors) > 0) == (len(test.expected) > 0) {
							t.Errorf("expected (len(mt.errors) > 0) = %t, got %t", len(test.expected) > 0, len(mt.errors) > 0)
						}
					}

				})
			}
		})
	}
}

func primitives() []test {
	return []test{
		{1, 1, nil, false},
		{1, 2, []string{"expected 1, got 2"}, false},
		{uint8(1), uint8(1), nil, false},
		{uint8(1), uint8(2), []string{"expected 1, got 2"}, false},
		{1.5, 1.5, nil, false},
		{1.5, 2.5, []string{"expected 1.5, got 2.5"}, false},
		{"a", "a", nil, false},
		{"a", "b", []string{"expected a, got b"}, false},

		{"a", 1, []string{"expected a, got 1"}, false},
	}
}

func structs() []test {
	type s struct {
		i int
	}
	type anotherS struct {
		i int
	}

	type withPointer struct {
		i *int
	}

	return []test{
		{s{1}, s{1}, nil, false},
		{s{1}, s{2}, []string{`not equal
differences:
--- expected
+++ actual
@@ -1,3 +1,3 @@
 teq_test.s{
-  i: [int(1)],
+  i: [int(2)],
 }
`}, false},
		{s{1}, anotherS{1}, []string{"expected {1}, got {1}"}, false},

		{withPointer{ref(1)}, withPointer{ref(1)}, nil, false},
		{withPointer{ref(1)}, withPointer{ref(2)}, []string{"expected {1}, got {2}"}, true},
	}
}

func slices() []test {
	return []test{
		{[]int{1, 2}, []int{1, 2}, nil, false},
		{[]int{1, 2}, []int{2, 1}, []string{`not equal
differences:
--- expected
+++ actual
@@ -1,4 +1,4 @@
 [
+  int(2),
   int(1),
-  int(2),
 ]
`}, false},
	}
}

func maps() []test {
	return []test{
		{map[string]int{"a": 1}, map[string]int{"a": 1}, nil, false},
		{map[string]int{"a": 1}, map[string]int{"a": 2}, []string{"expected map[a:1], got map[a:2]"}, false},
		{map[string]int{"a": 1}, map[string]int{"b": 1}, []string{"expected map[a:1], got map[b:1]"}, false},
		{map[string]int{"a": 0}, map[string]int{}, []string{"expected map[a:0], got map[]"}, false},

		{
			map[int]map[string]int{
				1: {"abc": 1},
			},
			map[int]map[string]int{
				1: {"abc": 1},
			},
			nil,
			false,
		},
		{
			map[int]map[string]int{
				1: {"abc": 1},
			},
			map[int]map[string]int{
				1: {"abc": 2},
			},
			[]string{"expected map[1:map[abc:1]], got map[1:map[abc:2]]"},
			false,
		},
	}
}

func recursions() []test {
	type privateRecursiveStruct struct {
		i int
		r *privateRecursiveStruct
	}
	r1_1 := privateRecursiveStruct{1, nil}
	r1_1.r = &r1_1
	r1_2 := privateRecursiveStruct{1, nil}
	r1_2.r = &r1_2
	r1_3 := privateRecursiveStruct{2, nil}
	r1_3.r = &r1_3
	r1_4 := privateRecursiveStruct{4, nil}
	r1_5 := privateRecursiveStruct{4, nil}
	r1_4.r = &r1_5
	r1_5.r = &r1_4
	r1_6 := privateRecursiveStruct{4, nil}
	r1_6.r = &r1_6

	type PublicRecursiveStruct struct {
		I int
		R *PublicRecursiveStruct
	}
	r2_1 := PublicRecursiveStruct{1, nil}
	r2_1.R = &r2_1
	r2_2 := PublicRecursiveStruct{1, nil}
	r2_2.R = &r2_2

	var r3_1 []any
	r3_1 = append(r3_1, 1, 2, r3_1)
	var r3_2 []any
	r3_2 = append(r3_2, 1, 2, r3_2)

	return []test{
		{r1_1, r1_1, nil, false},
		{r1_1, r1_2, nil, false},
		{r1_1, r1_3, []string{"expected {1, <cyclic>}, got {2, <cyclic>}"}, true},
		{r1_4, r1_5, nil, false},
		{r1_4, r1_6, nil, false},

		{r2_1, r2_1, nil, false},
		{r2_1, r2_2, nil, false},

		{r3_1, r3_1, nil, false},
		{r3_1, r3_2, nil, false},
	}
}

func ref[T any](v T) *T {
	return &v
}
