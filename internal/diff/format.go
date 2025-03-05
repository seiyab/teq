package diff

import (
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type format1 struct {
	value    reflect.Value
	original diffTree
	format   func(reflect.Value) string
}

func (f format1) docs() []doc.Doc {
	if f.format != nil && f.value.CanInterface() {
		return []doc.Doc{
			doc.Inline(printCustom(f.format, f.value)),
		}
	}
	return f.original.docs()
}

func (f format1) loss() float64 {
	return f.original.loss()
}

type format2 struct {
	left     reflect.Value
	right    reflect.Value
	original diffTree
	format   func(reflect.Value) string
}

func (m format2) docs() []doc.Doc {
	if m.format != nil && m.left.CanInterface() && m.right.CanInterface() {
		l := printCustom(m.format, m.left)
		r := printCustom(m.format, m.right)
		return []doc.Doc{
			doc.Inline(l).Left(),
			doc.Inline(r).Right(),
		}
	}

	l := printStringer(m.left)
	r := printStringer(m.right)
	if l == nil || r == nil {
		return m.original.docs() // fallback
	}
	return []doc.Doc{
		doc.Inline(*l).Left(),
		doc.Inline(*r).Right(),
	}
}

func (m format2) loss() float64 {
	return m.original.loss()
}
