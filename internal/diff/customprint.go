package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

var stringerType = reflect.TypeFor[fmt.Stringer]()

func printStringer(v reflect.Value) doc.Doc {
	if !v.Type().Implements(stringerType) || !v.CanInterface() {
		return nil
	}
	for v.Kind() == reflect.Interface && !v.IsNil() && v.Type() == stringerType {
		v = v.Elem()
	}

	m, ok := v.Interface().(fmt.Stringer)
	if !ok {
		return nil
	}
	b := m.String()
	return doc.BothInline(quote(string(b))).
		AddPrefix(fmt.Sprintf("%s(", v.Type().String())).
		AddSuffix(")")
}

func printCustom(f func(reflect.Value) string, v reflect.Value) doc.Doc {
	return doc.BothInline(quote(f(v))).
		AddPrefix(fmt.Sprintf("%s(", v.Type().String())).
		AddSuffix(")")
}
