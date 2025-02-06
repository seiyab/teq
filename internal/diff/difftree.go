package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type DiffTree struct {
	loss    float64
	split   bool
	entries []entry
	left    reflect.Value
	right   reflect.Value
}

func same(v reflect.Value) DiffTree {
	return DiffTree{
		loss:    0,
		entries: entreisOf(v),
		left:    v,
		right:   v,
	}
}

func imbalanced(v reflect.Value) DiffTree {
	return same(v) // smells bad :(
}

func eachSide(left, right reflect.Value) DiffTree {
	return DiffTree{
		loss:  1,
		split: true,
		entries: []entry{
			{value: imbalanced(left), leftOnly: true},
			{value: imbalanced(right), rightOnly: true},
		},
		left:  left,
		right: right,
	}
}

func (d DiffTree) Format() string {
	o := d.docs()
	return doc.PrintDoc(o)
}

func (d DiffTree) docs() []doc.Doc {
	if d.split {
		if len(d.entries) != 2 {
			panic("unexpected entries length")
		}
		var ds []doc.Doc
		l := d.entries[0].value.docs()
		for _, d := range l {
			ds = append(ds, d.Left())
		}
		r := d.entries[1].value.docs()
		for _, d := range r {
			ds = append(ds, d.Right())
		}
		return ds
	}
	if d.left.Kind() != d.right.Kind() {
		panic("kind mismatch: shouldn't happen. it's a bug if you see this")
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

func lossForEntries(es []entry) float64 {
	if len(es) == 0 {
		return 0
	}
	const max = 0.9
	total := 0.
	for _, e := range es {
		if e.leftOnly || e.rightOnly {
			total += 1
			continue
		}
		total += e.value.loss
	}
	return total / (max * float64(len(es)))
}
