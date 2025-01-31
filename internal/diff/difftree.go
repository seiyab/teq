package diff

import (
	"reflect"
	"strings"
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
	if d.left.Kind() != d.right.Kind() {
		panic("kind mismatch: not implemented")
	}
	switch d.left.Kind() {
	case reflect.String:
		return strings.Join([]string{
			"- " + d.left.String(),
			"+ " + d.right.String(),
		}, "\n")
	}
	panic("kind not implemented")
}

type entry struct {
	key   string
	value DiffTree
	left  bool
	right bool
}
