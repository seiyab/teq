package teq_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/seiyab/teq"
)

type test struct {
	a        any
	b        any
	expected []string
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
		{"interfaces", interfaces()},
		{"channels", channels()},
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

					for i, e := range test.expected {
						if mt.errors[i] != e {
							t.Errorf("expected %q, got %q at i = %d", e, mt.errors[i], i)
						}
						assert.Equal(t, e, mt.errors[i])
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
		{1, 1, nil},
		{1, 2, []string{"expected 1, got 2"}},
		{uint8(1), uint8(1), nil},
		{uint8(1), uint8(2), []string{"expected 1, got 2"}},
		{1.5, 1.5, nil},
		{1.5, 2.5, []string{"expected 1.5, got 2.5"}},
		{"a", "a", nil},
		{"a", "b", []string{"expected a, got b"}},

		{"a", 1, []string{"expected a, got 1"}},
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
		{s{1}, s{1}, nil},
		{s{1}, s{2}, []string{`not equal
differences:
--- expected
+++ actual
  teq_test.s{
-   i: 1,
+   i: 2,
  }`}},
		{s{1}, anotherS{1}, []string{"expected {1}, got {1}"}},

		{withPointer{ref(1)}, withPointer{ref(1)}, nil},
		{withPointer{ref(1)}, withPointer{ref(2)}, []string{`not equal
differences:
--- expected
+++ actual
  teq_test.withPointer{
-   i: &1,
+   i: &2,
  }`}},
	}
}

func slices() []test {
	return []test{
		{[]int{1, 2}, []int{1, 2}, nil},
		{[]int{1, 2}, []int{2, 1}, []string{`not equal
differences:
--- expected
+++ actual
  []int{
-   1,
    2,
+   1,
  }`}},
		{io.Reader(bytes.NewBuffer([]byte("a"))), io.Reader(bytes.NewBuffer(nil)), []string{
			`not equal
differences:
--- expected
+++ actual
- *bytes.Buffer("a")
+ *bytes.Buffer("")`,
		}},
	}
}

func maps() []test {
	return []test{
		{map[string]int{"a": 1}, map[string]int{"a": 1}, nil},
		{map[string]int{"a": 1}, map[string]int{"a": 2}, []string{`not equal
differences:
--- expected
+++ actual
  map[string]int{
-   "a": 1,
+   "a": 2,
  }`}},
		{map[string]int{"a": 1}, map[string]int{"b": 1}, []string{`not equal
differences:
--- expected
+++ actual
  map[string]int{
-   "a": 1,
+   "b": 1,
  }`}},
		{map[string]int{"a": 0}, map[string]int{}, []string{`not equal
differences:
--- expected
+++ actual
  map[string]int{
-   "a": 0,
  }`}},

		{
			map[int]map[string]int{
				1: {"abc": 1},
			},
			map[int]map[string]int{
				1: {"abc": 1},
			},
			nil,
		},
		{
			map[int]map[string]int{
				1: {"abc": 1},
			},
			map[int]map[string]int{
				1: {"abc": 2},
			},
			[]string{`not equal
differences:
--- expected
+++ actual
  map[int]map[string]int{
    1: map[string]int{
-     "abc": 1,
+     "abc": 2,
    },
  }`},
		},
		{
			map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
				"d": "4",
				"e": "5",
			},
			map[string]string{
				"a": "1",
				"b": "2",
				"c": "10000",
				"d": "4",
				"e": "5",
			},
			[]string{`not equal
differences:
--- expected
+++ actual
  map[string]string{
    "a": "1",
    "b": "2",
-   "c": "3",
+   "c": "10000",
    "d": "4",
    "e": "5",
  }`},
		},
	}
}

func interfaces() []test {
	return []test{
		{
			[]io.Reader{io.Reader(bytes.NewBuffer([]byte("a")))},
			[]io.Reader{io.Reader(bytes.NewBuffer([]byte("a")))},
			nil,
		},
		{
			[]io.Reader{
				bytes.NewBuffer([]byte("a")),
				bytes.NewBuffer([]byte("b")),
			},
			[]io.Reader{nil},
			[]string{`not equal
differences:
--- expected
+++ actual
  []io.Reader{
-   *bytes.Buffer("a"),
-   *bytes.Buffer("b"),
+   io.Reader(nil),
  }`}},
	}
}

func channels() []test {
	c1 := make(chan int)
	c2 := make(chan int)
	return []test{
		{c1, c1, nil},
		{c1, c2, []string{fmt.Sprintf("expected %p, got %p", c1, c2)}},
		{[]chan int{c1}, []chan int{c1}, nil},
		{[]chan int{c1}, []chan int{c2}, []string{`not equal
differences:
--- expected
+++ actual
  []chan int{
-   chan int,
+   chan int,
  }`}},
		{[]chan int{c1}, []chan int{nil}, []string{`not equal
differences:
--- expected
+++ actual
  []chan int{
-   chan int,
+   chan int(nil),
  }`}},
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
		{r1_1, r1_1, nil},
		{r1_1, r1_2, nil},
		{r1_1, r1_3, []string{`not equal
differences:
--- expected
+++ actual
  teq_test.privateRecursiveStruct{
-   i: 1,
+   i: 2,
    r: &teq_test.privateRecursiveStruct{
-     i: 1,
+     i: 2,
-     r: &<circular reference>,
+     r: &<circular reference>,
    },
  }`}},
		{r1_4, r1_5, nil},
		{r1_4, r1_6, nil},

		{r2_1, r2_1, nil},
		{r2_1, r2_2, nil},

		{r3_1, r3_1, nil},
		{r3_1, r3_2, nil},
	}
}

func ref[T any](v T) *T {
	return &v
}
