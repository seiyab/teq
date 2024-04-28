package teq

import (
	"reflect"
)

type Teq struct {
	MaxDepth int

	transforms map[reflect.Type]func(reflect.Value) reflect.Value
}

func New() Teq {
	return Teq{
		MaxDepth: 1_000,

		transforms: make(map[reflect.Type]func(reflect.Value) reflect.Value),
	}
}

func (teq Teq) Equal(t TestingT, expected, actual any) bool {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic in github.com/seiyab/teq. please report issue. message: %v", r)
		}
	}()
	ok := teq.equal(expected, actual)
	if !ok {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	return ok
}

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
		make(map[reflect.Type]bool),
		0,
	)
}
