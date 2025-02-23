package diff

import (
	"reflect"

	"github.com/seiyab/teq/internal/doc"
)

func eachSide(left, right reflect.Value) diffTree {
	return split{
		left:  pure(left),
		right: pure(right),
	}
}

type split struct {
	left  diffTree
	right diffTree
}

func (s split) docs() []doc.Doc {
	var ds []doc.Doc
	l := s.left.docs()
	for _, d := range l {
		ds = append(ds, d.Left())
	}
	r := s.right.docs()
	for _, d := range r {
		ds = append(ds, d.Right())
	}
	return ds
}

func (s split) loss() float64 {
	return 1
}
