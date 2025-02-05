package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type printNext func(t DiffTree) []doc.Doc
type printFunc = func(t DiffTree, n printNext) []doc.Doc

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

func notImplementedPrint(t DiffTree, n printNext) []doc.Doc {
	panic("not implemented")
}

func printSlice(t DiffTree, nx printNext) []doc.Doc {
	if t.loss == 0 {
		return []doc.Doc{
			doc.BothInline(t.left.Type().String() + "{ ... }"),
		}
	}
	var items []doc.Doc
	for _, e := range t.entries {
		docs := nx(e.value)
		for _, d := range docs {
			if e.leftOnly {
				d = d.Left()
			} else if e.rightOnly {
				d = d.Right()
			}
			items = append(items, d.AddSuffix(","))
		}
	}

	return []doc.Doc{
		doc.Block(
			doc.BothInline(t.left.Type().String()+"{"),
			items,
			doc.BothInline("}"),
		),
	}
}

func printString(t DiffTree, _ printNext) []doc.Doc {
	if t.loss == 0 {
		return []doc.Doc{
			doc.BothInline(quote(t.left.String())),
		}
	}
	return []doc.Doc{
		doc.LeftInline(quote(t.left.String())),
		doc.RightInline(quote(t.right.String())),
	}
}

func printStruct(t DiffTree, nx printNext) []doc.Doc {
	if t.loss == 0 {
		return []doc.Doc{
			doc.BothInline(t.left.Type().String() + "{ ... }"),
		}
	}
	var items []doc.Doc
	for _, e := range t.entries {
		items = append(items, printStructEntry(e, nx)...)
	}
	return []doc.Doc{
		doc.Block(
			doc.BothInline(t.left.Type().String()+"{"),
			items,
			doc.BothInline("}"),
		),
	}
}

func printStructEntry(e entry, nx printNext) []doc.Doc {
	docs := nx(e.value)
	var items []doc.Doc
	for _, d := range docs {
		items = append(
			items,
			d.AddPrefix(e.key+": ").AddSuffix(","),
		)
	}
	return items
}

var printInt = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%d", v.Int()) })
var printUint = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%d", v.Uint()) })
var printBool = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%t", v.Bool()) })
var printFloat = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%f", v.Float()) })
var printComplex = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%f", v.Complex()) })

func printPrimitive(f func(v reflect.Value) string) printFunc {
	return func(t DiffTree, _ printNext) []doc.Doc {
		if t.loss == 0 {
			return []doc.Doc{
				doc.BothInline(f(t.left)),
			}
		}
		return []doc.Doc{
			doc.LeftInline(f(t.left)),
			doc.RightInline(f(t.right)),
		}
	}
}
