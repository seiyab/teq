package diff

import (
	"fmt"
	"reflect"
)

type printNext func(t DiffTree) lines
type printFunc = func(t DiffTree, n printNext) lines

var printFuncs = map[reflect.Kind]printFunc{
	reflect.Array:      notImplementedPrint,
	reflect.Slice:      notImplementedPrint,
	reflect.Chan:       notImplementedPrint,
	reflect.Interface:  notImplementedPrint,
	reflect.Pointer:    notImplementedPrint,
	reflect.Struct:     notImplementedPrint,
	reflect.Map:        notImplementedPrint,
	reflect.Func:       notImplementedPrint,
	reflect.Int:        printInt,
	reflect.Int8:       printInt,
	reflect.Int16:      printInt,
	reflect.Int32:      printInt,
	reflect.Int64:      printInt,
	reflect.Uint:       printUint,
	reflect.Uint8:      printUint,
	reflect.Uint16:     printUint,
	reflect.Uint32:     printUint,
	reflect.Uint64:     printUint,
	reflect.Uintptr:    notImplementedPrint,
	reflect.String:     printString,
	reflect.Bool:       printBool,
	reflect.Float32:    printFloat,
	reflect.Float64:    printFloat,
	reflect.Complex64:  printComplex,
	reflect.Complex128: printComplex,
}

func notImplementedPrint(t DiffTree, n printNext) lines {
	panic("not implemented")
}

func printString(t DiffTree, _ printNext) lines {
	if t.loss == 0 {
		return lines{
			bothLine(quote(t.left.String())),
		}
	}
	return lines{
		leftLine(quote(t.left.String())),
		rightLine(quote(t.right.String())),
	}
}

var printInt = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%d", v.Int()) })
var printUint = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%d", v.Uint()) })
var printBool = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%t", v.Bool()) })
var printFloat = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%f", v.Float()) })
var printComplex = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%f", v.Complex()) })

func printPrimitive(f func(v reflect.Value) string) printFunc {
	return func(t DiffTree, _ printNext) lines {
		if t.loss == 0 {
			return lines{
				bothLine(f(t.left)),
			}
		}
		return lines{
			leftLine(f(t.left)),
			rightLine(f(t.right)),
		}
	}
}
