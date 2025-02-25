package diff

import (
	"fmt"

	"github.com/seiyab/teq/internal/doc"
)

type DiffTree struct {
	inner diffTree
}

type diffTree interface {
	docs() []doc.Doc
	loss() float64
}

var _ diffTree = split{}
var _ diffTree = mixed{}
var _ diffTree = cycle{}
var _ diffTree = nilNode{}
var _ diffTree = format1{}
var _ diffTree = format2{}
var _ diffTree = fail{}

func (d DiffTree) Format() string {
	o := d.inner.docs()
	return doc.PrintDoc(o)
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}
