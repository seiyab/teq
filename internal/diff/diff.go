package diff

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Differ struct {
	reflectEqual func(v1, v2 reflect.Value) bool
}

func New() Differ {
	return Differ{}
}

func (d Differ) Diff(x, y any) (DiffTree, error) {
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	return d.diff(v1, v2, make(map[visit]bool), 0)
}

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

const maxDepth = 100

func (d Differ) diff(
	v1, v2 reflect.Value,
	visited map[visit]bool,
	depth int,
) (DiffTree, error) {
	if depth > maxDepth {
		return DiffTree{}, fmt.Errorf("maximum depth exceeded")
	}
	if d.reflectEqual != nil {
		if d.reflectEqual(v1, v2) {
			return same(v1), nil
		}
	} else {
		if reflect.DeepEqual(v1.Interface(), v2.Interface()) {
			return same(v1), nil
		}
	}
	if !v1.IsValid() || !v2.IsValid() {
		return DiffTree{}, fmt.Errorf("not implemented")
	}
	if v1.Type() != v2.Type() {
		return DiffTree{}, fmt.Errorf("not implemented")
	}

	if hard(v1.Kind()) {
		if v1.CanAddr() && v2.CanAddr() {
			addr1 := v1.Addr().UnsafePointer()
			addr2 := v2.Addr().UnsafePointer()

			// Short circuit
			if uintptr(addr1) == uintptr(addr2) {
				return same(v1), nil
			}
			if uintptr(addr1) > uintptr(addr2) {
				// Canonicalize order to reduce number of entries in visited.
				addr1, addr2 = addr2, addr1
			}

			// Short circuit if references are already seen.
			typ := v1.Type()
			v := visit{addr1, addr2, typ}
			if visited[v] {
				return same(v1), nil
			}

			// Remember for later.
			visited[v] = true
		}
	}

	diffFunc, ok := diffFuncs[v1.Kind()]
	if !ok {
		panic("diff is not defined for " + v1.Type().String())
	}
	var n next = func(v1, v2 reflect.Value) (DiffTree, error) {
		return d.diff(v1, v2, visited, depth+1)
	}
	return diffFunc(v1, v2, n)
}

func hard(k reflect.Kind) bool {
	switch k {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return true
	}
	return false
}
