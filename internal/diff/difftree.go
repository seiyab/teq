package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type DiffTree struct {
	inner diffTree
}

type diffTree interface {
	docs() []doc.Doc
	loss() float64
}

var _ diffTree = split{}
var _ diffTree = mixed{}

type split struct {
	left  diffTree
	right diffTree
}

func (s split) docs() []doc.Doc {
	var ds []doc.Doc
	l := s.left.docs()
	for _, d := range l {
		ds = append(ds, d.Left())
	}
	r := s.right.docs()
	for _, d := range r {
		ds = append(ds, d.Right())
	}
	return ds
}

func (s split) loss() float64 {
	return 1
}

type mixed struct {
	distance float64
	sample   reflect.Value
	entries  []entry
}

func (m mixed) docs() []doc.Doc {
	f, ok := printFuncs[m.sample.Kind()]
	if !ok {
		panic("not implemented: " + m.sample.Kind().String())
	}
	return f(m)
}

func (m mixed) loss() float64 {
	return m.distance
}

func pure(v reflect.Value) diffTree {
	return mixed{
		distance: 0,
		sample:   v,
		entries:  entriesOf(v),
	}
}

func same(v reflect.Value) diffTree {
	return pure(v)
}

func imbalanced(v reflect.Value) diffTree {
	return pure(v)
}

func eachSide(left, right reflect.Value) diffTree {
	return split{
		left:  pure(left),
		right: pure(right),
	}
}

func (d DiffTree) Format() string {
	o := d.inner.docs()
	return doc.PrintDoc(o)
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

type entry struct {
	key       string
	value     diffTree
	leftOnly  bool
	rightOnly bool
}

func lossForKeyedEntries(es []entry) float64 {
	if len(es) == 0 {
		return 0.1
	}
	const max = 0.9
	total := 0.
	for _, e := range es {
		if e.leftOnly || e.rightOnly {
			total += 1
			continue
		}
		total += e.value.loss()
	}
	return total / (max * float64(len(es)))
}

func lossForIndexedEntries(es []entry) float64 {
	const max = 0.9
	if len(es) == 0 {
		return max
	}
	n := 0.
	total := 0.
	for _, e := range es {
		switch t := e.value.(type) {
		case split:
			total += 2
			n += 2
		case mixed:
			if e.leftOnly || e.rightOnly {
				n += 1
				total += 1
				break
			}
			total += t.loss()
			n += 1
		default:
			panic("unexpected type: " + fmt.Sprintf("%T", t))
		}
	}
	return max * total / n
}
