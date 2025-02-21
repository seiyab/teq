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
var _ diffTree = cycle{}
var _ diffTree = nilNode{}

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

func pure(v reflect.Value) diffTree {
	return recPure(v, make(map[visit]bool))
}

func recPure(val reflect.Value, visited map[visit]bool) diffTree {
	k := val.Kind()
	f, ok := entriesFuncs[k]
	if !ok {
		return mixed{
			distance: 0,
			sample:   val,
			entries:  nil,
		}
	}
	if !hard(val) || !val.CanAddr() {
		return mixed{
			distance: 0,
			sample:   val,
			entries: f(val, func(v reflect.Value) diffTree {
				return recPure(v, visited)
			}),
		}
	}

	vis := visit{ptr: val.Addr().UnsafePointer(), typ: val.Type()}
	if visited[vis] {
		return cycle{}
	}

	visited = cloneVisits(visited)
	visited[vis] = true

	return mixed{
		distance: 0,
		sample:   val,
		entries: f(val, func(v reflect.Value) diffTree {
			return recPure(v, visited)
		}),
	}
}

func same(v reflect.Value) diffTree {
	return pure(v)
}

func imbalanced(v reflect.Value) diffTree {
	return pure(v)
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

func eachSide(left, right reflect.Value) diffTree {
	return split{
		left:  pure(left),
		right: pure(right),
	}
}

type cycle struct{}

func (c cycle) docs() []doc.Doc {
	return []doc.Doc{
		doc.BothInline("<circular reference>"),
	}
}

func (c cycle) loss() float64 {
	return 0
}

func null(v reflect.Value) diffTree {
	return nilNode{v.Type()}
}

type nilNode struct{ ty reflect.Type }

func (n nilNode) docs() []doc.Doc {
	return []doc.Doc{
		doc.BothInline(fmt.Sprintf("%s(nil)", n.ty.String())),
	}
}

func (nilNode) loss() float64 {
	return 0
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
		case mixed, cycle, nilNode:
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
