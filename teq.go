package teq

import (
	"reflect"
)

type Teq struct{}

func (teq Teq) Equal(t TestingT, expected, actual any) bool {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic: %v", r)
		}
	}()
	ok := teq.equal(expected, actual)
	if !ok {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	return ok
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
	return teq.deepValueEqual(v1, v2, make(map[visit]bool), 0)
}
