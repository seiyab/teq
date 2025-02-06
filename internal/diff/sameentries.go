package diff

import "reflect"

type entriesFunc func(v reflect.Value, n nextEntries) []entry
type nextEntries func(v reflect.Value) []entry

func entriesOf(v reflect.Value) []entry {
	f, ok := entriesFuncs[v.Kind()]
	if !ok {
		return nil
	}
	return f(v, entriesOf)
}

var entriesFuncs = map[reflect.Kind]entriesFunc{}

func init() {
	entriesFuncs[reflect.Array] = entriesNotImplemented
	entriesFuncs[reflect.Slice] = sliceEntries
	entriesFuncs[reflect.Interface] = entriesNotImplemented
	entriesFuncs[reflect.Struct] = structEntries
	entriesFuncs[reflect.Map] = entriesNotImplemented
	entriesFuncs[reflect.String] = stringEntries
}

func entriesNotImplemented(v reflect.Value, n nextEntries) []entry {
	panic("not implemented")
}

func sliceEntries(v reflect.Value, nx nextEntries) []entry {
	var es []entry
	for i := 0; i < v.Len(); i++ {
		x := v.Index(i)
		es = append(es, entry{
			value: same(x),
		})
	}
	return es
}

func structEntries(v reflect.Value, nx nextEntries) []entry {
	var es []entry
	n := v.NumField()
	for i := 0; i < n; i++ {
		k := v.Type().Field(i).Name
		x := v.Field(i)
		es = append(es, entry{
			key:   k,
			value: same(x),
		})
	}
	return es
}

func stringEntries(v reflect.Value, nx nextEntries) []entry {
	return nil
}
