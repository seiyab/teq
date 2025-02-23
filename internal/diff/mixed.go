package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

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
	if m.sample.Type().Implements(textMarshalerType) && m.distance == 0 {
		mt := printMarshalText(m.sample)
		if mt != nil {
			return []doc.Doc{
				mt,
			}
		}
	}

	f, ok := printFuncs[m.sample.Kind()]
	if !ok {
		panic("not implemented: " + m.sample.Kind().String())
	}
	return f(m)
}

func (m mixed) loss() float64 {
	return m.distance
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
