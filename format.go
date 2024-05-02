package teq

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
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
	if k != reflect.Struct &&
		k != reflect.Map &&
		k != reflect.Slice &&
		k != reflect.Array &&
		k != reflect.String {
		return simple
	}
	if k == reflect.String && len(ve.String()) < 10 && len(va.String()) < 10 {
		return simple
	}

	diff := difflib.UnifiedDiff{
		A:        addLineBreak(teq.format(ve, 0)),
		B:        addLineBreak(teq.format(va, 0)),
		FromFile: "expected",
		ToFile:   "actual",
		Context:  1,
	}
	diffTxt, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return fmt.Sprintf("failed to get diff: %v", err)
	}
	if diffTxt == "" {
		return simple
	}
	return strings.Join([]string{
		"not equal",
		"differences:",
		diffTxt,
	}, "\n")
}

func (teq Teq) format(v reflect.Value, depth int) []string {
	if depth > teq.MaxDepth {
		return []string{"<max depth exceeded>"}
	}
	if !v.IsValid() {
		return []string{"<invalid>"}
	}

	fmtFn, ok := fmts[v.Kind()]
	if !ok {
		fmtFn = todoFmt
	}
	next := func(v reflect.Value) []string {
		return teq.format(v, depth+1)
	}
	return fmtFn(v, next)
}

var fmts = map[reflect.Kind]func(reflect.Value, func(reflect.Value) []string) []string{
	reflect.Array:      todoFmt,
	reflect.Slice:      sliceFmt,
	reflect.Interface:  todoFmt,
	reflect.Pointer:    todoFmt,
	reflect.Struct:     structFmt,
	reflect.Map:        todoFmt,
	reflect.Func:       todoFmt,
	reflect.Int:        intFmt,
	reflect.Int8:       intFmt,
	reflect.Int16:      intFmt,
	reflect.Int32:      intFmt,
	reflect.Int64:      intFmt,
	reflect.Uint:       uintFmt,
	reflect.Uint8:      uintFmt,
	reflect.Uint16:     uintFmt,
	reflect.Uint32:     uintFmt,
	reflect.Uint64:     uintFmt,
	reflect.Uintptr:    uintFmt,
	reflect.String:     stringFmt,
	reflect.Bool:       boolFmt,
	reflect.Float32:    floatFmt,
	reflect.Float64:    floatFmt,
	reflect.Complex64:  complexFmt,
	reflect.Complex128: complexFmt,
}

func todoFmt(v reflect.Value, next func(reflect.Value) []string) []string {
	return []string{fmt.Sprintf("<%s>", v.String())}
}

func sliceFmt(v reflect.Value, next func(reflect.Value) []string) []string {
	if v.Len() == 0 {
		return []string{"[]"}
	}
	result := make([]string, 0, v.Len()+2)
	result = append(result, "[")
	for i := 0; i < v.Len(); i++ {
		result = append(result, indent(fmt.Sprintf("%s,", next(v.Index(i))[0])))
	}
	result = append(result, "]")
	return result
}

func structFmt(v reflect.Value, next func(reflect.Value) []string) []string {
	open := fmt.Sprintf("%s{", v.Type())
	close := "}"
	if v.NumField() == 0 {
		return []string{open + close}
	}
	result := make([]string, 0, v.NumField()+2)
	result = append(result, open)
	for i := 0; i < v.NumField(); i++ {
		result = append(result, indent(fmt.Sprintf("%s: %s,", v.Type().Field(i).Name, next(v.Field(i)))))
	}
	result = append(result, close)
	return result
}

func intFmt(v reflect.Value, _ func(reflect.Value) []string) []string {
	return []string{fmt.Sprintf("%s(%d)", v.Type(), v.Int())}
}
func uintFmt(v reflect.Value, _ func(reflect.Value) []string) []string {
	return []string{fmt.Sprintf("%s(%d)", v.Type(), v.Uint())}
}
func stringFmt(v reflect.Value, _ func(reflect.Value) []string) []string {
	return []string{v.String()}
}
func boolFmt(v reflect.Value, _ func(reflect.Value) []string) []string {
	return []string{fmt.Sprintf("%t", v.Bool())}
}
func floatFmt(v reflect.Value, _ func(reflect.Value) []string) []string {
	return []string{fmt.Sprintf("%s(%f)", v.Type(), v.Float())}
}
func complexFmt(v reflect.Value, _ func(reflect.Value) []string) []string {
	return []string{fmt.Sprintf("%s(%f, %f)", v.Type(), real(v.Complex()), imag(v.Complex()))}
}

func indent(s string) string {
	return fmt.Sprintf("  %s", s)
}
func addLineBreak(ss []string) []string {
	result := make([]string, 0, len(ss))
	for _, s := range ss {
		result = append(result, fmt.Sprintf("%s\n", s))
	}
	return result
}
