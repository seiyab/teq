package diff

import (
	"fmt"
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

func null(v reflect.Value) diffTree {
	return nilNode{v.Type()}
}

type nilNode struct{ ty reflect.Type }

func (n nilNode) docs() []doc.Doc {
	return []doc.Doc{
		doc.Inline(fmt.Sprintf("%s(nil)", n.ty.String())),
	}
}

func (nilNode) loss() float64 {
	return 0
}
