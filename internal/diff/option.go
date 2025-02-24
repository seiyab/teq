package diff

import (
	"reflect"
)

// Option is an option for DiffString.
type Option func(*differ) *differ

// WithReflectEqual sets a function to compare reflect.Value.
func WithReflectEqual(fn func(v1, v2 reflect.Value) bool) Option {
	return func(d *differ) *differ {
		d.reflectEqual = fn
		return d
	}
}

// WithFormat sets a function to format a value.
func WithFormat[T any](fn func(T) string) Option {
	return func(d *differ) *differ {
		var zero T
		if d.formats == nil {
			d.formats = make(map[reflect.Type]func(reflect.Value) string)
		}
		d.formats[reflect.TypeOf(zero)] = func(v reflect.Value) string {
			r := reflect.ValueOf(fn)
			return r.Call([]reflect.Value{v})[0].String()
		}
		return d
	}
}
