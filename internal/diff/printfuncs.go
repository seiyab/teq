package diff

import (
	"fmt"
	"reflect"
)

type printNext func(t DiffTree) lines
type printFunc = func(t DiffTree, n printNext) lines

var printFuncs = map[reflect.Kind]printFunc{
	reflect.Array:      notImplementedPrint,
	reflect.Slice:      printSlice,
	reflect.Chan:       notImplementedPrint,
	reflect.Interface:  notImplementedPrint,
	reflect.Pointer:    notImplementedPrint,
	reflect.Struct:     printStruct,
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

func printSlice(t DiffTree, nx printNext) lines {
	if t.loss == 0 {
		return lines{
			bothLine(t.left.Type().String() + "{ ... }"),
		}
	}
	var result lines
	result.add(bothLine(t.left.Type().String() + "{").open())
	for _, e := range t.entries {
		ls := nx(e.value)
		if e.leftOnly {
			ls.left()
		} else if e.rightOnly {
			ls.right()
		}
		result.concat(ls)
	}
	result.add(bothLine("}").close())

	return result
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

func printStruct(t DiffTree, nx printNext) lines {
	if t.loss == 0 {
		return lines{
			bothLine(t.left.Type().String() + "{ ... }"),
		}
	}
	var result lines
	result.add(bothLine(t.left.Type().String() + "{").open())
	for _, e := range t.entries {
		result.concat(printStructEntry(e, nx))
	}
	result.add(bothLine("}").close())

	return result
}

func printStructEntry(e entry, nx printNext) lines {
	ls := nx(e.value)
	if len(ls) == 0 {
		panic("unexpected empty lines")
	}
	z := len(ls) - 1
	if ls[0].onLeft && ls[0].onRight {
		ls[0] = ls[0].overrideText(e.key + ": " + ls[0].text)
		if ls[z].onLeft && ls[z].onRight {
			ls[z] = ls[z].overrideText(ls[z].text + ",")
		} else {
			ls.add(bothLine(","))
		}
	} else if len(ls) == 2 {
		ls[0] = ls[0].overrideText(e.key + ": " + ls[0].text + ",")
		ls[1] = ls[1].overrideText(e.key + ": " + ls[1].text + ",")
	} else {
		ls = append(lines{leftLine(e.key + ":").open()}, ls...)
		ls.add(leftLine(",").close())
	}
	return ls
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
