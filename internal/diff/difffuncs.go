package diff

import (
	"fmt"
	"reflect"
)

type next func(v1, v2 reflect.Value) (DiffTree, error)
type diffFunc = func(v1, v2 reflect.Value, n next) (DiffTree, error)

var diffFuncs = map[reflect.Kind]diffFunc{
	reflect.Array:      notImplemented,
	reflect.Slice:      notImplemented,
	reflect.Chan:       notImplemented,
	reflect.Interface:  notImplemented,
	reflect.Pointer:    notImplemented,
	reflect.Struct:     notImplemented,
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

func notImplemented(v1, v2 reflect.Value, _ next) (DiffTree, error) {
	return DiffTree{}, fmt.Errorf("not implemented")
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
		return DiffTree{loss: 1, left: v1, right: v2}, nil
	}
}
