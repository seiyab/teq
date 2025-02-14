package diff

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type next func(v1, v2 reflect.Value) (diffTree, error)
type diffFunc = func(v1, v2 reflect.Value, n next) (diffTree, error)

var diffFuncs = map[reflect.Kind]diffFunc{
	reflect.Array:      sliceDiff,
	reflect.Slice:      sliceDiff,
	reflect.Chan:       notImplemented,
	reflect.Interface:  notImplemented,
	reflect.Pointer:    notImplemented,
	reflect.Struct:     structDiff,
	reflect.Map:        mapDiff,
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

func notImplemented(v1, v2 reflect.Value, n next) (diffTree, error) {
	return nil, fmt.Errorf("not implemented")
}

func sliceDiff(v1, v2 reflect.Value, nx next) (diffTree, error) {
	if v1.Type() != v2.Type() {
		return nil, fmt.Errorf("unexpected type mismatch")
	}
	if v1.Kind() == reflect.Slice && (v1.IsNil() || v2.IsNil()) {
		if v1.IsNil() && v2.IsNil() {
			return same(v1), nil
		}
		return eachSide(v1, v2), nil
	}
	es, err := sliceMixedEntries(v1, v2, nx)
	if err != nil {
		return nil, err
	}

	return mixed{
		distance: lossForIndexedEntries(es),
		entries:  es,
		sample:   v1,
	}, nil
}

func structDiff(v1, v2 reflect.Value, nx next) (diffTree, error) {
	if v1.Type() != v2.Type() {
		return eachSide(v1, v2), nil
	}
	entries := make([]entry, 0, v1.NumField())
	for i, n := 0, v1.NumField(); i < n; i++ {
		key := v1.Type().Field(i).Name
		vd, err := nx(field(v1, i), field(v2, i))
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry{key: key, value: vd, leftOnly: true, rightOnly: true})
	}
	return mixed{
		distance: lossForKeyedEntries(entries),
		entries:  entries,
		sample:   v1,
	}, nil
}

func stringDiff(v1, v2 reflect.Value, _ next) (diffTree, error) {
	s1, s2 := v1.String(), v2.String()
	if s1 == s2 {
		return same(v1), nil
	}
	lines1 := strings.Split(s1, "\n")
	lines2 := strings.Split(s2, "\n")
	if len(lines1) == 1 || len(lines2) == 1 {
		return eachSide(v1, v2), nil
	}
	es, err := multiLineStringEntries(lines1, lines2)
	if err != nil {
		return nil, err
	}
	return mixed{
		distance: lossForIndexedEntries(es),
		entries:  es,
		sample:   v1,
	}, nil
}

var intDiff = primitiveDiff(func(v reflect.Value) int64 { return v.Int() })
var uintDiff = primitiveDiff(func(v reflect.Value) uint64 { return v.Uint() })
var boolDiff = primitiveDiff(func(v reflect.Value) bool { return v.Bool() })
var floatDiff = primitiveDiff(func(v reflect.Value) float64 { return v.Float() })
var complexDiff = primitiveDiff(func(v reflect.Value) complex128 { return v.Complex() })

func primitiveDiff[T comparable](f func(v reflect.Value) T) diffFunc {
	return func(v1, v2 reflect.Value, _ next) (diffTree, error) {
		if f(v1) == f(v2) {
			return same(v1), nil
		}
		return eachSide(v1, v2), nil
	}
}

func mapDiff(v1, v2 reflect.Value, nx next) (diffTree, error) {
	if v1.Type() != v2.Type() {
		return eachSide(v1, v2), nil
	}
	if v1.IsNil() || v2.IsNil() {
		if v1.IsNil() && v2.IsNil() {
			return same(v1), nil
		}
		return eachSide(v1, v2), nil
	}

	var keys []reflect.Value
	iter1 := v1.MapRange()
	for iter1.Next() {
		keys = append(keys, iter1.Key())
	}
	iter2 := v2.MapRange()
	for iter2.Next() {
		k := iter2.Key()
		if v1.MapIndex(k).IsValid() {
			continue
		}
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return compareMapKey(keys[i], keys[j])
	})

	var entries []entry
	for _, k := range keys {
		val1 := v1.MapIndex(k)
		val2 := v2.MapIndex(k)

		if !val1.IsValid() && !val2.IsValid() {
			continue // shouldn't happen
		}

		if !val1.IsValid() {
			entries = append(entries, entry{
				key:       stringifyKey(k),
				value:     same(val2),
				rightOnly: true,
			})
			continue
		}

		if !val2.IsValid() {
			entries = append(entries, entry{
				key:      stringifyKey(k),
				value:    same(val1),
				leftOnly: true,
			})
			continue
		}

		d, err := nx(val1, val2)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry{
			key:   stringifyKey(k),
			value: d,
		})
	}

	return mixed{
		distance: lossForKeyedEntries(entries),
		entries:  entries,
		sample:   v1,
	}, nil
}

func compareMapKey(a, b reflect.Value) bool {
	if a.Kind() != b.Kind() {
		return stringifyKey(a) < stringifyKey(b)
	}
	switch a.Kind() {
	case reflect.String:
		return a.String() < b.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() < b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.Bool:
		return !a.Bool() && b.Bool()
	default:
		return stringifyKey(a) < stringifyKey(b)
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
