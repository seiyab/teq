package diff

import (
	"fmt"
	"reflect"
	"unsafe"
)

func DiffString(x, y any) (string, error) {
	d := New()
	t, err := d.Diff(x, y)
	if err != nil {
		return "", err
	}
	return t.Format(), nil
}

type Differ struct {
	reflectEqual func(v1, v2 reflect.Value) bool
}

func New() Differ {
	return Differ{}
}

func (d Differ) Diff(x, y any) (DiffTree, error) {
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	p := diffProcess{differ: d}
	t, err := p.diff(v1, v2, 0)
	if err != nil {
		return DiffTree{}, err
	}
	return DiffTree{inner: t}, nil
}

type diffProcess struct {
	differ       Differ
	leftVisited  map[visit]bool
	rightVisited map[visit]bool
}

type visit struct {
	ptr unsafe.Pointer
	typ reflect.Type
}

const maxDepth = 500

func (p diffProcess) diff(
	v1, v2 reflect.Value,
	depth int,
) (diffTree, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("maximum depth exceeded")
	}

	d := p.differ
	if d.reflectEqual != nil {
		if d.reflectEqual(v1, v2) {
			return same(v1), nil
		}
	} else if lightDeepEqual(v1, v2) {
		return same(v1), nil
	}
	if !v1.IsValid() || !v2.IsValid() {
		return nil, fmt.Errorf("invalid value")
	}
	if v1.Type() != v2.Type() {
		return eachSide(v1, v2), nil
	}
	if d.reflectEqual == nil {
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
	var n next = func(v1, v2 reflect.Value) (diffTree, error) {
		return p.diff(v1, v2, depth+1)
	}
	t, err := diffFunc(v1, v2, n)
	if err != nil {
		return nil, err
	}
	if v1.Type().Implements(textMarshalerType) {
		t = marshal{left: v1, right: v2, real: t}
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
	l := make(map[visit]bool, len(p.leftVisited))
	for k, b := range p.leftVisited {
		l[k] = b
	}
	r := make(map[visit]bool, len(p.rightVisited))
	for k, b := range p.rightVisited {
		r[k] = b
	}
	return diffProcess{
		differ:       p.differ,
		leftVisited:  l,
		rightVisited: r,
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
