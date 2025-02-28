package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

func (d diffProcess) pure(val reflect.Value) diffTree {
	fm := func(m mixed) diffTree {
		if f, ok := d.differ.formats[val.Type()]; ok {
			return format1{value: val, original: m, format: f}
		} else if val.Type().Implements(stringerType) && val.Kind() != reflect.String {
			return format1{value: val, original: m, format: nil}
		}
		return m
	}

	k := val.Kind()
	f, ok := entriesFuncs[k]
	if !ok {
		return fm(mixed{
			distance: 0,
			sample:   val,
			entries:  nil,
		})
	}
	if !hard(val) || !val.CanAddr() {
		return fm(mixed{
			distance: 0,
			sample:   val,
			entries:  f(val, d),
		})
	}

	vis := visit{ptr: val.Addr().UnsafePointer(), typ: val.Type()}
	if d.pureVisited[vis] {
		return cycle{}
	}

	d = d.clone()
	d.pureVisited[vis] = true

	return fm(mixed{
		distance: 0,
		sample:   val,
		entries:  f(val, d),
	})
}

type mixed struct {
	distance float64
	sample   reflect.Value
	entries  []entry
}

func (m mixed) docs() []doc.Doc {
	if m.sample.Type().Implements(stringerType) && m.distance == 0 {
		mt := printStringer(m.sample)
		if mt != nil {
			return []doc.Doc{
				doc.Inline(*mt),
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
	keyName   string
	keyValue  reflect.Value
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
		case mixed, cycle, nilNode, format1, format2:
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
