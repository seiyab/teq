package teq

import (
	"reflect"
)

// Teq is a object for deep equality comparison.
type Teq struct {
	// MaxDepth is the maximum depth of the comparison. Default is 1000.
	MaxDepth int

	transforms map[reflect.Type]func(reflect.Value) reflect.Value
	formats    map[reflect.Type]func(reflect.Value) string
	equals     map[reflect.Type]func(reflect.Value, reflect.Value) bool
}

// New returns new instance of Teq.
func New() Teq {
	return Teq{
		MaxDepth: 1_000,

		transforms: make(map[reflect.Type]func(reflect.Value) reflect.Value),
		formats:    make(map[reflect.Type]func(reflect.Value) string),
		equals:     make(map[reflect.Type]func(reflect.Value, reflect.Value) bool),
	}
}

// Equal perform deep equality check and report error if not equal.
func (teq Teq) Equal(t TestingT, expected, actual any) bool {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic in github.com/seiyab/teq. please report issue. message: %v", r)
		}
	}()
	ok := teq.equal(expected, actual)
	if !ok {
		t.Error(teq.report(expected, actual))
	}
	return ok
}

// NotEqual perform deep equality check and report error if equal.
func (teq Teq) NotEqual(t TestingT, expected, actual any) bool {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic in github.com/seiyab/teq. please report issue. message: %v", r)
		}
	}()
	ok := !teq.equal(expected, actual)
	if !ok {
		if reflect.DeepEqual(expected, actual) {
			t.Error("reflect.DeepEqual(expected, actual) == true.")
		} else {
			t.Errorf("expected %v != %v", expected, actual)
			t.Log("reflect.DeepEqual(expected, actual) == false. maybe transforms made them equal.")
		}
	}
	return ok

}

// AddTransform adds a transform function to Teq.
// The transform function must have only one argument and one return value.
// The argument type is the type to be transformed.
// If the passed transform function is not valid, it will panic.
// The transformed value will be used for equality check instead of the original value.
// The transformed value and its internal values won't be transformed to prevent infinite recursion.
func (teq *Teq) AddTransform(transform any) {
	ty := reflect.TypeOf(transform)
	if ty.Kind() != reflect.Func {
		panic("transform must be a function")
	}
	if ty.NumIn() != 1 {
		panic("transform must have only one argument")
	}
	if ty.NumOut() != 1 {
		panic("transform must have only one return value")
	}
	trValue := reflect.ValueOf(transform)
	reflectTransform := func(v reflect.Value) reflect.Value {
		return trValue.Call([]reflect.Value{v})[0]
	}
	teq.transforms[ty.In(0)] = reflectTransform
}

// AddFormat adds a format function to Teq.
// The format function must have only one argument and one return value.
// The argument type is the type to be formatted.
// If the passed format function is not valid, it will panic.
// The formatted string will be shown instead of the original value in the error report when the values are not equal.
func (teq *Teq) AddFormat(format any) {
	ty := reflect.TypeOf(format)
	if ty.Kind() != reflect.Func {
		panic("format must be a function")
	}
	if ty.NumIn() != 1 {
		panic("format must have only one argument")
	}
	if ty.NumOut() != 1 {
		panic("format must have only one return value")
	}
	if ty.Out(0).Kind() != reflect.String {
		panic("format must return string")
	}
	formatValue := reflect.ValueOf(format)
	reflectFormat := func(v reflect.Value) string {
		return formatValue.Call([]reflect.Value{v})[0].String()
	}
	teq.formats[ty.In(0)] = reflectFormat
}

// AddEqual adds an equal function to Teq.
// The equal function must have two arguments with the same type and one return value of bool.
// If the passed equal function is not valid, it will panic.
// The equal function will be used for equality check instead of the default equality check.
func (teq *Teq) AddEqual(equal any) {
	ty := reflect.TypeOf(equal)
	if ty.Kind() != reflect.Func {
		panic("equal must be a function")
	}
	if ty.NumIn() != 2 {
		panic("equal must have two arguments")
	}
	if ty.In(0) != ty.In(1) {
		panic("equal must have two arguments with the same type")
	}
	if ty.NumOut() != 1 {
		panic("equal must have only one return value")
	}
	if ty.Out(0).Kind() != reflect.Bool {
		panic("equal must return bool")
	}
	equalValue := reflect.ValueOf(equal)
	reflectEqual := func(v1, v2 reflect.Value) bool {
		return equalValue.Call([]reflect.Value{v1, v2})[0].Bool()
	}
	teq.equals[ty.In(0)] = reflectEqual
}

func (teq Teq) equal(x, y any) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return teq.deepValueEqual(
		v1, v2,
		make(map[visit]bool),
		0,
	)
}
