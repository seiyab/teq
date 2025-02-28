package diff

import (
	"fmt"
	"reflect"
)

var stringerType = reflect.TypeFor[fmt.Stringer]()

func printStringer(v reflect.Value) *string {
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
	s := fmt.Sprintf("%s(%q)", v.Type().String(), m.String())
	return &s
}

func printCustom(f func(reflect.Value) string, v reflect.Value) string {
	return fmt.Sprintf("%s(%q)", v.Type().String(), f(v))
}
