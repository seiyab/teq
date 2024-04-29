// Some code is written referencing following codes:
// - deepequal.go in "reflect" package authored by Go Authors
// - deepequal.go in "github.com/weaveworks/scope/test/reflect" package authored by Weaveworks Ltd

package teq

import (
	"reflect"
	"unsafe"
)

// During deepValueEqual, must keep track of checks that are
// in progress. The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func (teq Teq) deepValueEqual(
	v1, v2 reflect.Value,
	visited map[visit]bool,
	depth int,
) bool {
	if depth > teq.MaxDepth {
		panic("maximum depth exceeded")
	}
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	tr, ok := teq.transforms[v1.Type()]
	if ok {
		t1 := tr(v1)
		t2 := tr(v2)
		newTeq := New()
		newTeq.MaxDepth = teq.MaxDepth
		return newTeq.deepValueEqual(t1, t2, visited, depth)
	}

	if hard(v1.Kind()) {
		if v1.CanAddr() && v2.CanAddr() {
			addr1 := v1.Addr().UnsafePointer()
			addr2 := v2.Addr().UnsafePointer()

			// Short circuit
			if uintptr(addr1) == uintptr(addr2) {
				return true
			}
			if uintptr(addr1) > uintptr(addr2) {
				// Canonicalize order to reduce number of entries in visited.
				addr1, addr2 = addr2, addr1
			}

			// Short circuit if references are already seen.
			typ := v1.Type()
			v := visit{addr1, addr2, typ}
			if visited[v] {
				return true
			}

			// Remember for later.
			visited[v] = true
		}
	}

	eqFn, ok := eqs[v1.Kind()]
	if !ok {
		panic("not implemented")
	}
	var n next = func(v1, v2 reflect.Value) bool {
		return teq.deepValueEqual(v1, v2, visited, depth+1)
	}
	return eqFn(teq, v1, v2, n)
}

type next func(v1, v2 reflect.Value) bool

var eqs = map[reflect.Kind]func(teq Teq, v1, v2 reflect.Value, nx next) bool{
	reflect.Array:      arrayEq,
	reflect.Slice:      todo,
	reflect.Interface:  interfaceEq,
	reflect.Pointer:    pointerEq,
	reflect.Struct:     structEq,
	reflect.Map:        todo,
	reflect.Func:       todo,
	reflect.Int:        intEq,
	reflect.Int8:       intEq,
	reflect.Int16:      intEq,
	reflect.Int32:      intEq,
	reflect.Int64:      intEq,
	reflect.Uint:       uintEq,
	reflect.Uint8:      uintEq,
	reflect.Uint16:     uintEq,
	reflect.Uint32:     uintEq,
	reflect.Uint64:     uintEq,
	reflect.Uintptr:    uintEq,
	reflect.String:     stringEq,
	reflect.Bool:       boolEq,
	reflect.Float32:    floatEq,
	reflect.Float64:    floatEq,
	reflect.Complex64:  complexEq,
	reflect.Complex128: complexEq,
}

func hard(k reflect.Kind) bool {
	switch k {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return true
	}
	return false
}

func todo(teq Teq, v1, v2 reflect.Value, nx next) bool {
	panic("not implemented")
}

func arrayEq(teq Teq, v1, v2 reflect.Value, nx next) bool {
	for i := 0; i < v1.Len(); i++ {
		if !nx(v1.Index(i), v2.Index(i)) {
			return false
		}
	}
	return true
}

func interfaceEq(teq Teq, v1, v2 reflect.Value, nx next) bool {
	if v1.IsNil() || v2.IsNil() {
		return v1.IsNil() == v2.IsNil()
	}
	return nx(v1.Elem(), v2.Elem())
}

func pointerEq(teq Teq, v1, v2 reflect.Value, nx next) bool {
	if v1.UnsafePointer() == v2.UnsafePointer() {
		return true
	}
	return nx(v1.Elem(), v2.Elem())
}

func structEq(teq Teq, v1, v2 reflect.Value, nx next) bool {
	for i, n := 0, v1.NumField(); i < n; i++ {
		if !nx(field(v1, i), field(v2, i)) {
			return false
		}
	}
	return true
}

func intEq(_ Teq, v1, v2 reflect.Value, _ next) bool     { return v1.Int() == v2.Int() }
func uintEq(_ Teq, v1, v2 reflect.Value, _ next) bool    { return v1.Uint() == v2.Uint() }
func stringEq(_ Teq, v1, v2 reflect.Value, _ next) bool  { return v1.String() == v2.String() }
func boolEq(_ Teq, v1, v2 reflect.Value, _ next) bool    { return v1.Bool() == v2.Bool() }
func floatEq(_ Teq, v1, v2 reflect.Value, _ next) bool   { return v1.Float() == v2.Float() }
func complexEq(_ Teq, v1, v2 reflect.Value, _ next) bool { return v1.Complex() == v2.Complex() }
