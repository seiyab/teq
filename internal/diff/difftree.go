package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type DiffTree struct {
	loss    float64
	entries []entry
	left    reflect.Value
	right   reflect.Value
}

func same(v reflect.Value) DiffTree {
	return DiffTree{loss: 0, left: v, right: v}
}

func imbalanced(v reflect.Value) DiffTree {
	return DiffTree{loss: 0, left: v, right: v}
}

func (d DiffTree) Format() string {
	o := d.docs()
	return doc.PrintDoc(o)
}

func (d DiffTree) docs() []doc.Doc {
	if d.left.Kind() != d.right.Kind() {
		panic("kind mismatch: not implemented")
	}
	f, ok := printFuncs[d.left.Kind()]
	if !ok {
		panic("not implemented: " + d.left.Kind().String())
	}
	return f(d, func(t DiffTree) []doc.Doc {
		return t.docs()
	})
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

type entry struct {
	key       string
	value     DiffTree
	leftOnly  bool
	rightOnly bool
}
