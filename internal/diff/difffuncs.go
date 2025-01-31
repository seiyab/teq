package diff

import (
	"fmt"
	"reflect"
)

type next func(v1, v2 reflect.Value) (DiffTree, error)

var diffFuncs = map[reflect.Kind]func(reflect.Value, reflect.Value, next) (DiffTree, error){
	reflect.Array:      notImplemented,
	reflect.Slice:      notImplemented,
	reflect.Chan:       notImplemented,
	reflect.Interface:  notImplemented,
	reflect.Pointer:    notImplemented,
	reflect.Struct:     notImplemented,
	reflect.Map:        notImplemented,
	reflect.Func:       notImplemented,
	reflect.Int:        notImplemented,
	reflect.Int8:       notImplemented,
	reflect.Int16:      notImplemented,
	reflect.Int32:      notImplemented,
	reflect.Int64:      notImplemented,
	reflect.Uint:       notImplemented,
	reflect.Uint8:      notImplemented,
	reflect.Uint16:     notImplemented,
	reflect.Uint32:     notImplemented,
	reflect.Uint64:     notImplemented,
	reflect.Uintptr:    notImplemented,
	reflect.String:     stringDiff,
	reflect.Bool:       notImplemented,
	reflect.Float32:    notImplemented,
	reflect.Float64:    notImplemented,
	reflect.Complex64:  notImplemented,
	reflect.Complex128: notImplemented,
}

func notImplemented(v1, v2 reflect.Value, _ next) (DiffTree, error) {
	return DiffTree{}, fmt.Errorf("not implemented")
}

func stringDiff(v1, v2 reflect.Value, _ next) (DiffTree, error) {
	if v1.String() == v2.String() {
		return same(v1), nil
	}
	return DiffTree{
		loss:  1,
		left:  v1,
		right: v2,
	}, nil
}
