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

func (d DiffTree) Format() string {
	if d.loss == 0 {
		return ""
	}
	l := d.lines()
	return l.print()
}

func (d DiffTree) lines() lines {
	if d.loss == 0 {
		if isNaive(d.left.Kind()) {
			return lines{bothLine(formatNaive(d.left))}
		}
		panic("not implemented")
	}
	if d.left.Kind() != d.right.Kind() {
		panic("kind mismatch: not implemented")
	}
	if isNaive(d.left.Kind()) {
		return lines{
			leftLine(formatNaive(d.left)),
			rightLine(formatNaive(d.right)),
		}
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

func isNaive(k reflect.Kind) bool {
	switch k {
	case
		// reflect.String, // -- can be multiline
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		// reflect.Uintptr, // -- not sure
		reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

func formatNaive(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%f", v.Complex())
	}
	panic("unexpected kind" + v.Kind().String())
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
