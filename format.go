package teq

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/seiyab/akashi"
)

func (teq Teq) report(expected, actual any) string {
	simple := fmt.Sprintf("expected %v, got %v", expected, actual)
	if expected == nil || actual == nil {
		return simple
	}
	ve := reflect.ValueOf(expected)
	va := reflect.ValueOf(actual)
	if ve.Type() != va.Type() {
		return simple
	}
	k := ve.Kind()
	_, ok := teq.formats[ve.Type()]
	if !ok &&
		k != reflect.Struct &&
		k != reflect.Map &&
		k != reflect.Slice &&
		k != reflect.Array &&
		k != reflect.String &&
		k != reflect.Pointer {
		return simple
	}

	head := []string{
		"not equal",
		"differences:",
		"--- expected",
		"+++ actual",
	}
	options := []akashi.Option{}
	for _, f := range teq.formats {
		options = append(options, akashi.WithFormat(f))
	}
	diff := akashi.DiffString(expected, actual, options...)
	return strings.Join(append(head, diff), "\n")
}
