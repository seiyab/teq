package diff

import (
	"fmt"
	"reflect"
)

type DiffTree struct {
	loss    int
	entries []entry
	left    reflect.Value
	right   reflect.Value
}

func same(v reflect.Value) DiffTree {
	return DiffTree{loss: 0, left: v, right: v}
}

func (d DiffTree) Format() string {
	if d.loss == 0 {
		return ""
	}
	l := d.lines()
	return l.print()
}

func (d DiffTree) lines() lines {
	if d.loss == 0 {
		return lines{}
	}
	if d.left.Kind() != d.right.Kind() {
		panic("kind mismatch: not implemented")
	}
	switch d.left.Kind() {
	case reflect.String:
		return lines{
			leftLine(quote(d.left.String())),
			rightLine(quote(d.right.String())),
		}
	}
	panic("not implemented kind" + d.left.Kind().String())
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

type entry struct {
	key   string
	value DiffTree
	left  bool
	right bool
}
