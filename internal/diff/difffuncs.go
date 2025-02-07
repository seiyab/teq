package diff

import (
	"fmt"
	"reflect"
)

type next func(v1, v2 reflect.Value) (DiffTree, error)
type diffFunc = func(v1, v2 reflect.Value, n next) (DiffTree, error)

var diffFuncs = map[reflect.Kind]diffFunc{
	reflect.Array:      notImplemented,
	reflect.Slice:      sliceDiff,
	reflect.Chan:       notImplemented,
	reflect.Interface:  notImplemented,
	reflect.Pointer:    notImplemented,
	reflect.Struct:     structDiff,
	reflect.Map:        notImplemented,
	reflect.Func:       notImplemented,
	reflect.Int:        intDiff,
	reflect.Int8:       intDiff,
	reflect.Int16:      intDiff,
	reflect.Int32:      intDiff,
	reflect.Int64:      intDiff,
	reflect.Uint:       uintDiff,
	reflect.Uint8:      uintDiff,
	reflect.Uint16:     uintDiff,
	reflect.Uint32:     uintDiff,
	reflect.Uint64:     uintDiff,
	reflect.Uintptr:    uintDiff,
	reflect.String:     stringDiff,
	reflect.Bool:       boolDiff,
	reflect.Float32:    floatDiff,
	reflect.Float64:    floatDiff,
	reflect.Complex64:  complexDiff,
	reflect.Complex128: complexDiff,
}

func notImplemented(v1, v2 reflect.Value, n next) (DiffTree, error) {
	return DiffTree{}, fmt.Errorf("not implemented")
}

func sliceDiff(v1, v2 reflect.Value, nx next) (DiffTree, error) {
	if v1.Type() != v2.Type() {
		return DiffTree{}, fmt.Errorf("unexpected type mismatch")
	}
	if v1.IsNil() || v2.IsNil() {
		if v1.IsNil() && v2.IsNil() {
			return same(v1), nil
		}
		return eachSide(v1, v2), nil
	}
	es, err := sliceMixedEntries(v1, v2, nx)
	if err != nil {
		return DiffTree{}, err
	}

	return DiffTree{
		loss:    lossForIndexedEntries(es),
		entries: es,
		left:    v1,
		right:   v2,
	}, nil
}

func structDiff(v1, v2 reflect.Value, nx next) (DiffTree, error) {
	if v1.Type() != v2.Type() {
		return eachSide(v1, v2), nil
	}
	entries := make([]entry, 0, v1.NumField())
	for i, n := 0, v1.NumField(); i < n; i++ {
		key := v1.Type().Field(i).Name
		vd, err := nx(field(v1, i), field(v2, i))
		if err != nil {
			return DiffTree{}, err
		}
		entries = append(entries, entry{key: key, value: vd, leftOnly: true, rightOnly: true})
	}
	return DiffTree{
		loss:    lossForKeyedEntries(entries),
		entries: entries,
		left:    v1,
		right:   v2,
	}, nil
}

var intDiff = primitiveDiff(func(v reflect.Value) int64 { return v.Int() })
var uintDiff = primitiveDiff(func(v reflect.Value) uint64 { return v.Uint() })
var stringDiff = primitiveDiff(func(v reflect.Value) string { return v.String() })
var boolDiff = primitiveDiff(func(v reflect.Value) bool { return v.Bool() })
var floatDiff = primitiveDiff(func(v reflect.Value) float64 { return v.Float() })
var complexDiff = primitiveDiff(func(v reflect.Value) complex128 { return v.Complex() })

func primitiveDiff[T comparable](f func(v reflect.Value) T) diffFunc {
	return func(v1, v2 reflect.Value, _ next) (DiffTree, error) {
		if f(v1) == f(v2) {
			return same(v1), nil
		}
		return eachSide(v1, v2), nil
	}
}

func field(v reflect.Value, idx int) reflect.Value {
	f1 := v.Field(idx)
	if f1.CanAddr() {
		return f1
	}
	vc := reflect.New(v.Type()).Elem()
	vc.Set(v)
	rf := vc.Field(idx)
	return reflect.NewAt(rf.Type(), rf.Addr().UnsafePointer()).Elem()
}
