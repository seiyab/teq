package teq

import (
	"fmt"
	"reflect"
	"sort"
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
		k != reflect.String &&
		k != reflect.Pointer {
		return simple
	}
	if k == reflect.String {
		if len(ve.String()) < 10 && len(va.String()) < 10 {
			return simple
		}
		if strings.Contains(ve.String(), "\n") || strings.Contains(va.String(), "\n") {
			r, ok := richReport(
				difflib.SplitLines(ve.String()),
				difflib.SplitLines(va.String()),
			)
			if !ok {
				return simple
			}
			return r
		}
	}

	r, ok := richReport(
		teq.format(ve, 0).diffSequence(),
		teq.format(va, 0).diffSequence(),
	)
	if !ok {
		return simple
	}
	return r
}

func richReport(a []string, b []string) (string, bool) {
	diff := difflib.UnifiedDiff{
		A:        a,
		B:        b,
		FromFile: "expected",
		ToFile:   "actual",
		Context:  1,
	}
	diffTxt, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return fmt.Sprintf("failed to get diff: %v", err), false
	}
	if diffTxt == "" {
		return "", false
	}
	return strings.Join([]string{
		"not equal",
		"differences:",
		diffTxt,
	}, "\n"), true
}

func (teq Teq) format(v reflect.Value, depth int) lines {
	if depth > teq.MaxDepth {
		return linesOf("<max depth exceeded>")
	}
	if !v.IsValid() {
		return linesOf("<invalid>")
	}

	ty := v.Type()
	if fm, ok := teq.formats[ty]; ok {
		return linesOf(fm(v))
	}

	fmtFn, ok := fmts[v.Kind()]
	if !ok {
		fmtFn = todoFmt
	}
	next := func(v reflect.Value) lines {
		return teq.format(v, depth+1)
	}
	return fmtFn(v, next)
}

var fmts = map[reflect.Kind]func(reflect.Value, func(reflect.Value) lines) lines{
	reflect.Array:      arrayFmt,
	reflect.Slice:      sliceFmt,
	reflect.Interface:  interfaceFmt,
	reflect.Pointer:    pointerFmt,
	reflect.Struct:     structFmt,
	reflect.Map:        mapFmt,
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

func todoFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("<%s>", v.String()))
}

func arrayFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	open := fmt.Sprintf("%s{", v.Type().String())
	close := "}"
	if v.Len() == 0 {
		return linesOf(open + close)
	}
	result := make(lines, 0, v.Len()+2)
	result = append(result, lineOf(open))
	for i := 0; i < v.Len(); i++ {
		elem := next(v.Index(i)).followedBy(",")
		result = append(result, elem.indent()...)
	}
	result = append(result, lineOf(close))
	return result

}

func sliceFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	open := fmt.Sprintf("[]%s{", v.Type().Elem().String())
	close := "}"
	if v.Len() == 0 {
		return linesOf(open + close)
	}
	result := make(lines, 0, v.Len()+2)
	result = append(result, lineOf(open))
	for i := 0; i < v.Len(); i++ {
		elem := next(v.Index(i)).followedBy(",")
		result = append(result, elem.indent()...)
	}
	result = append(result, lineOf(close))
	return result
}

func interfaceFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	open := fmt.Sprintf("%s(", v.Type().String())
	close := ")"
	if v.IsNil() {
		return linesOf(open + "<nil>" + close)
	}
	return next(v.Elem()).ledBy(open).followedBy(close)
}

func pointerFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	if v.IsNil() {
		return linesOf("<nil>")
	}
	return next(v.Elem()).ledBy("*")
}

func structFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	open := fmt.Sprintf("%s{", v.Type().String())
	close := "}"
	if v.NumField() == 0 {
		return linesOf(open + close)
	}
	result := make(lines, 0, v.NumField()+2)
	result = append(result, lineOf(open))
	for i := 0; i < v.NumField(); i++ {
		entry := next(v.Field(i)).
			ledBy(v.Type().Field(i).Name + ": ").
			followedBy(",")
		result = append(result, entry.indent()...)
	}
	result = append(result, lineOf(close))
	return result
}

func mapFmt(v reflect.Value, next func(reflect.Value) lines) lines {
	open := fmt.Sprintf("map[%s]%s{", v.Type().Key(), v.Type().Elem())
	close := "}"
	if v.Len() == 0 {
		return linesOf(open + close)
	}
	result := make(lines, 0, v.Len()+2)
	result = append(result, lineOf(open))

	type entry struct {
		key   string
		lines lines
	}
	entries := make([]entry, 0, v.Len())
	for _, key := range v.MapKeys() {
		var e entry
		keyLines := next(key)
		e.key = keyLines.key()
		valLines := next(v.MapIndex(key))
		e.lines = keyValue(keyLines, valLines)
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})
	for _, e := range entries {
		result = append(result, e.lines.indent()...)
	}
	result = append(result, lineOf(close))
	return result
}

func intFmt(v reflect.Value, _ func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("%s(%d)", v.Type(), v.Int()))
}
func uintFmt(v reflect.Value, _ func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("%s(%d)", v.Type(), v.Uint()))
}
func stringFmt(v reflect.Value, _ func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("%q", v.String()))
}
func boolFmt(v reflect.Value, _ func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("%t", v.Bool()))
}
func floatFmt(v reflect.Value, _ func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("%s(%f)", v.Type(), v.Float()))
}
func complexFmt(v reflect.Value, _ func(reflect.Value) lines) lines {
	return linesOf(fmt.Sprintf("%s(%f, %f)", v.Type(), real(v.Complex()), imag(v.Complex())))
}
