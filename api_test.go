package teq_test

import (
	"fmt"
	"testing"

	"github.com/seiyab/teq"
)

type test struct {
	a        any
	b        any
	expected []string
}

func TestEqual(t *testing.T) {
	assert := teq.Teq{}

	tests := primitives()
	for _, test := range tests {
		name := fmt.Sprintf("%T(%v) == %T(%v)", test.a, test.a, test.b, test.b)
		t.Run(name, func(t *testing.T) {
			mt := &mockT{}
			assert.Equal(mt, test.a, test.b)
			if len(mt.errors) != len(test.expected) {
				t.Fatalf("expected %d errors, got %d", len(test.expected), len(mt.errors))
			}
			for i, e := range test.expected {
				if mt.errors[i] != e {
					t.Errorf("expected %v, got %v at i = %d", e, mt.errors[i], i)
				}
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
		{"a", "a", nil},
		{"a", "b", []string{"expected a, got b"}},

		{"a", 1, []string{"expected a, got 1"}},
	}
}
