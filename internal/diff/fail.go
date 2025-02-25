package diff

import (
	"fmt"

	"github.com/seiyab/teq/internal/doc"
)

type fail struct {
	difference float64
	message    string
}

func (f fail) docs() []doc.Doc {
	return []doc.Doc{
		doc.Inline(fmt.Sprintf("<failed to print diff: %q>", f.message)),
	}
}

func (f fail) loss() float64 {
	return f.difference
}
