package teq

import (
	"reflect"
)

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
