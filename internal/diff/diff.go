package diff

import (
	"fmt"
	"reflect"
	"unsafe"
)

// DiffString returns a string that represents the difference between x and y.
func DiffString(x, y any, options ...Option) (string, error) {
	d := differ{}
	for _, opt := range options {
		d = *opt(&d)
	}
	t, err := d.diff(x, y)
	if err != nil {
		return "", err
	}
	return t.Format(), nil
}

type differ struct {
	reflectEqual func(v1, v2 reflect.Value) bool
	formats      formats
}

type formats map[reflect.Type]func(reflect.Value) string

func (d differ) diff(x, y any) (DiffTree, error) {
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	p := diffProcess{differ: d}
	t, err := p.diff(v1, v2)
	if err != nil {
		return DiffTree{}, err
	}
	return DiffTree{inner: t}, nil
}

type diffProcess struct {
	differ       differ
	depth        int
	leftVisited  map[visit]bool
	rightVisited map[visit]bool
	pureVisited  map[visit]bool
}

type visit struct {
	ptr unsafe.Pointer
	typ reflect.Type
}

const maxDepth = 500

func (p diffProcess) diff(
	v1, v2 reflect.Value,
) (diffTree, error) {
	if p.depth > maxDepth {
		return nil, fmt.Errorf("maximum depth exceeded")
	}
	p.depth = p.depth + 1

	d := p.differ
	if d.reflectEqual != nil {
		if d.reflectEqual(v1, v2) {
			return p.pure(v1), nil
		}
	} else if lightDeepEqual(v1, v2) {
		return p.pure(v1), nil
	}
	if !v1.IsValid() || !v2.IsValid() {
		return nil, fmt.Errorf("invalid value")
	}
	if v1.Type() != v2.Type() {
		return p.eachSide(v1, v2), nil
	}

	p, cyclic := p.cycle(v1, v2)
	if cyclic {
		return split{
			left:  cycle{},
			right: cycle{},
		}, nil
	}

	diffFunc, ok := diffFuncs[v1.Kind()]
	if !ok {
		panic("diff is not defined for " + v1.Type().String())
	}
	t, err := diffFunc(v1, v2, p)
	if err != nil {
		return nil, err
	}
	if f, ok := d.formats[v1.Type()]; ok {
		t = format2{left: v1, right: v2, original: t, format: f}
	} else if v1.Type().Implements(stringerType) {
		t = format2{left: v1, right: v2, original: t}
	}
	return t, nil
}

func lightDeepEqual(v1 reflect.Value, v2 reflect.Value) bool {
	if v1.Type() != v2.Type() {
		return false
	}
	if v1.CanInterface() && v2.CanInterface() {
		return reflect.DeepEqual(v1.Interface(), v2.Interface())
	}
	if v1.CanAddr() && v2.CanAddr() && v1.Addr().Pointer() == v2.Addr().Pointer() {
		return true
	}
	return false // can't go better until go 1.20
}

func (p diffProcess) cycle(v1 reflect.Value, v2 reflect.Value) (diffProcess, bool) {
	if !hard(v1) && !hard(v2) {
		return p, false
	}
	leftCycle := false
	rightCycle := false
	p = p.clone()
	if hard(v1) && v1.CanAddr() {
		addr := v1.Addr().UnsafePointer()
		vis := visit{ptr: addr, typ: v1.Type()}
		leftCycle = p.leftVisited[vis]
		p.leftVisited[vis] = true
	}
	if hard(v2) && v2.CanAddr() {
		addr := v2.Addr().UnsafePointer()
		vis := visit{ptr: addr, typ: v2.Type()}
		rightCycle = p.rightVisited[vis]
		p.rightVisited[vis] = true
	}

	return p, leftCycle && rightCycle
}

func (p diffProcess) clone() diffProcess {
	return diffProcess{
		differ:       p.differ,
		leftVisited:  cloneVisits(p.leftVisited),
		rightVisited: cloneVisits(p.rightVisited),
		pureVisited:  cloneVisits(p.pureVisited),
	}
}

func hard(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map:
		return !v.IsNil()
	case reflect.Struct, reflect.Array:
		return true
	}
	return false
}
