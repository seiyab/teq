package teq

import "reflect"

func copyTransformed(m map[reflect.Type]bool) map[reflect.Type]bool {
	c := make(map[reflect.Type]bool, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}
