package diff

import (
	"fmt"
	"reflect"
)

type entriesFunc func(v reflect.Value, d diffProcess) []entry

func cloneVisits(vis map[visit]bool) map[visit]bool {
	v := make(map[visit]bool, len(vis))
	for k := range vis {
		v[k] = true
	}
	return v
}

var entriesFuncs = map[reflect.Kind]entriesFunc{}

func init() {
	entriesFuncs[reflect.Array] = sliceEntries
	entriesFuncs[reflect.Slice] = sliceEntries
	entriesFuncs[reflect.Interface] = interfaceEntries
	entriesFuncs[reflect.Pointer] = pointerEntries
	entriesFuncs[reflect.Struct] = structEntries
	entriesFuncs[reflect.Map] = mapEntries
	entriesFuncs[reflect.String] = stringEntries
}

func sliceEntries(v reflect.Value, d diffProcess) []entry {
	var es []entry
	for i := 0; i < v.Len(); i++ {
		x := v.Index(i)
		es = append(es, entry{
			value: d.pure(x),
		})
	}
	return es
}

func interfaceEntries(v reflect.Value, d diffProcess) []entry {
	if v.IsNil() {
		return nil
	}
	return []entry{
		{value: d.pure(v.Elem())},
	}
}

func pointerEntries(v reflect.Value, d diffProcess) []entry {
	if v.IsNil() {
		return nil
	}
	return []entry{
		{value: d.pure(v.Elem())},
	}
}

func structEntries(v reflect.Value, d diffProcess) []entry {
	var es []entry
	n := v.NumField()
	for i := 0; i < n; i++ {
		k := v.Type().Field(i).Name
		x := v.Field(i)
		es = append(es, entry{
			key:   k,
			value: d.pure(x),
		})
	}
	return es
}

func stringEntries(v reflect.Value, d diffProcess) []entry {
	return nil
}

func mapEntries(v reflect.Value, d diffProcess) []entry {
	var es []entry
	iter := v.MapRange()
	for iter.Next() {
		k := iter.Key()
		x := iter.Value()
		es = append(es, entry{
			key:   stringifyKey(k),
			value: d.pure(x),
		})
	}
	return es
}

func stringifyKey(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%g", v.Complex())
	case reflect.Interface:
		if v.IsNil() {
			return fmt.Sprintf("%s(<nil>)", v.Type().String())
		}
		return stringifyKey(v.Elem())
	default:
		return fmt.Sprint(v.Interface())
	}
}
