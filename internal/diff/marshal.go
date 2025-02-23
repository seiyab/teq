package diff

import (
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

type marshal struct {
	left  reflect.Value
	right reflect.Value
	real  diffTree
}

func (m marshal) docs() []doc.Doc {
	l := printMarshalText(m.left)
	r := printMarshalText(m.right)
	if l == nil || r == nil {
		return m.real.docs() // fallback
	}
	return []doc.Doc{
		l.Left(),
		r.Right(),
	}
}

func (m marshal) loss() float64 {
	return m.real.loss()
}
