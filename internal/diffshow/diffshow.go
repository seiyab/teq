package diffshow

import "reflect"

type segment string

func Compare(a any, b any, equal func(any, any) bool) Diff {
	if a == nil || b == nil {
		panic("TODO: not implemented")
	}
	return compare(reflect.ValueOf(a), reflect.ValueOf(b), nil)
}

func compare(a reflect.Value, b reflect.Value, path []segment) Diff {
	/*
		if a == nil || b == nil {
			diff := Diff{path: path}
			if a != nil {
				diff.a = ref(reflect.ValueOf(a))
			}
			if b != nil {
				diff.b = ref(reflect.ValueOf(b))
			}
			return diff
		}
	*/
	if a.Type() != b.Type() {
		return Diff{
			path: path,
			a:    a,
			b:    b,
		}
	}

	diffFn, ok := diffs[a.Kind()]
	if !ok {
		panic("unsupported kind: " + a.Kind().String())
	}
	var n next = func(v1, v2 reflect.Value, path []segment) Diff {
		return compare(a, b, path)
	}
	return diffFn(a, b, path, n)
}

var diffs = map[reflect.Kind]func(v1, v2 reflect.Value, path []segment, nx next) Diff{
	reflect.Array:      nil,
	reflect.Slice:      nil,
	reflect.Chan:       nil,
	reflect.Interface:  nil,
	reflect.Pointer:    nil,
	reflect.Struct:     nil,
	reflect.Map:        nil,
	reflect.Func:       nil,
	reflect.Int:        primitive,
	reflect.Int8:       primitive,
	reflect.Int16:      primitive,
	reflect.Int32:      primitive,
	reflect.Int64:      primitive,
	reflect.Uint:       primitive,
	reflect.Uint8:      primitive,
	reflect.Uint16:     primitive,
	reflect.Uint32:     primitive,
	reflect.Uint64:     primitive,
	reflect.Uintptr:    nil,
	reflect.String:     primitive,
	reflect.Bool:       primitive,
	reflect.Float32:    primitive,
	reflect.Float64:    primitive,
	reflect.Complex64:  nil,
	reflect.Complex128: nil,
}

type next func(a, b reflect.Value, path []segment) Diff

func primitive(a, b reflect.Value, path []segment, _ next) Diff {
	return Diff{
		path: path,
		a:    a,
		b:    b,
	}
}
