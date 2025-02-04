package diff

import (
	"fmt"
	"reflect"
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
	return DiffTree{loss: 1, left: v, right: v}
}

func (d DiffTree) Format() string {
	if d.loss == 0 {
		return ""
	}
	l := d.lines()
	return l.print()
}

func (d DiffTree) lines() lines {
	if d.left.Kind() != d.right.Kind() {
		panic("kind mismatch: not implemented")
	}
	f, ok := printFuncs[d.left.Kind()]
	if !ok {
		panic("not implemented: " + d.left.Kind().String())
	}
	return f(d, func(t DiffTree) lines {
		return t.lines()
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
