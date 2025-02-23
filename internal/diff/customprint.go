package diff

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

var textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()

func printMarshalText(v reflect.Value) doc.Doc {
	if !v.Type().Implements(textMarshalerType) {
		return nil
	}
	m := v.Interface().(encoding.TextMarshaler)
	b, err := m.MarshalText()
	if err != nil {
		return nil
	}
	return doc.BothInline(quote(string(b))).
		AddPrefix(fmt.Sprintf("%s(", v.Type().String())).
		AddSuffix(")")
}
