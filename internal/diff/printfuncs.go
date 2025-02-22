package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type printFunc = func(t mixed) []doc.Doc

var printFuncs = map[reflect.Kind]printFunc{
	reflect.Array:      printSlice,
	reflect.Slice:      printSlice,
	reflect.Chan:       printChan,
	reflect.Interface:  printInterface,
	reflect.Pointer:    printPointer,
	reflect.Struct:     printStruct,
	reflect.Map:        printMap,
	reflect.Func:       printFn,
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
	reflect.Uintptr:    printUintptr,
	reflect.String:     printString,
	reflect.Bool:       printBool,
	reflect.Float32:    printFloat,
	reflect.Float64:    printFloat,
	reflect.Complex64:  printComplex,
	reflect.Complex128: printComplex,
}

func printSlice(m mixed) []doc.Doc {
	var items []doc.Doc
	for _, e := range m.entries {
		docs := e.value.docs()
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
			doc.BothInline(m.sample.Type().String()+"{"),
			items,
			doc.BothInline("}"),
		),
	}
}

func printChan(m mixed) []doc.Doc {
	ty := m.sample.Type().String()
	if m.sample.IsNil() {
		return []doc.Doc{
			doc.BothInline(fmt.Sprintf("%s(nil)", ty)),
		}
	}
	return []doc.Doc{
		doc.BothInline(ty),
	}
}

func printString(m mixed) []doc.Doc {
	if m.loss() == 0 {
		return []doc.Doc{
			doc.BothInline(quote(m.sample.String())),
		}
	}
	var items []doc.Doc
	for _, e := range m.entries {
		t, ok := e.value.(mixed)
		if !ok {
			panic("unexpected type")
		}
		s := t.sample.String()
		if e.leftOnly {
			items = append(items, doc.LeftInline(s))
		} else if e.rightOnly {
			items = append(items, doc.RightInline(s))
		} else {
			items = append(items, doc.BothInline(s))
		}
	}
	return []doc.Doc{
		doc.Block(
			doc.BothInline(m.sample.Type().Name()+"("),
			items,
			doc.BothInline(")"),
		),
	}
}

func printInterface(m mixed) []doc.Doc {
	if m.sample.IsNil() {
		return []doc.Doc{
			doc.BothInline(fmt.Sprintf("%s(nil)", m.sample.Type().String())),
		}
	}
	return m.entries[0].value.docs()
}

func printPointer(m mixed) []doc.Doc {
	if m.sample.IsNil() {
		return []doc.Doc{
			doc.BothInline(fmt.Sprintf("%s(nil)", m.sample.Type().String())),
		}
	}
	docs := m.entries[0].value.docs()
	for i := range docs {
		docs[i] = docs[i].AddPrefix("&")
	}
	return docs
}

func printStruct(m mixed) []doc.Doc {
	var items []doc.Doc
	for _, e := range m.entries {
		items = append(items, printStructEntry(e)...)
	}
	return []doc.Doc{
		doc.Block(
			doc.BothInline(m.sample.Type().String()+"{"),
			items,
			doc.BothInline("}"),
		),
	}
}

func printStructEntry(e entry) []doc.Doc {
	docs := e.value.docs()
	var items []doc.Doc
	for _, d := range docs {
		items = append(
			items,
			d.AddPrefix(e.key+": ").AddSuffix(","),
		)
	}
	return items
}

func printMap(m mixed) []doc.Doc {
	if m.sample.IsNil() {
		return []doc.Doc{
			doc.BothInline(fmt.Sprintf("%s(nil)", m.sample.Type().String())),
		}
	}

	var items []doc.Doc
	for _, e := range m.entries {
		docs := e.value.docs()
		for _, d := range docs {
			if e.leftOnly {
				d = d.Left()
			} else if e.rightOnly {
				d = d.Right()
			}
			items = append(items, d.AddPrefix(e.key+": ").AddSuffix(","))
		}
	}

	return []doc.Doc{
		doc.Block(
			doc.BothInline(m.sample.Type().String()+"{"),
			items,
			doc.BothInline("}"),
		),
	}
}

func printFn(m mixed) []doc.Doc {
	ty := m.sample.Type().String()
	if m.sample.IsNil() {
		return []doc.Doc{
			doc.BothInline(fmt.Sprintf("%s(nil)", ty)),
		}
	}
	return []doc.Doc{
		doc.BothInline(fmt.Sprintf("%s { ... }", ty)),
	}
}

var printInt = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%d", v.Int()) })
var printUint = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%d", v.Uint()) })
var printUintptr = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%s(%d)", v.Type().String(), v.Uint()) })
var printBool = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%t", v.Bool()) })
var printFloat = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%f", v.Float()) })
var printComplex = printPrimitive(func(v reflect.Value) string { return fmt.Sprintf("%f", v.Complex()) })

func printPrimitive(f func(v reflect.Value) string) printFunc {
	return func(m mixed) []doc.Doc {
		return []doc.Doc{doc.BothInline(f(m.sample))}
	}
}
