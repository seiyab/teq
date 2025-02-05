package doc

import "strings"

const indentSize = 2

type Doc interface {
	print(depth int, buf *buffer)
	Left() Doc
	Right() Doc
	AddPrefix(prefix string) Doc
	AddSuffix(suffix string) Doc
}

var _ Doc = block{}
var _ Doc = inline{}

func PrintDoc(ds []Doc) string {
	var b buffer
	for _, d := range ds {
		d.print(0, &b)
	}
	return strings.Join(b.lines, "\n")
}

type buffer struct {
	lines []string
}

func (b *buffer) push(s string) {
	b.lines = append(b.lines, s)
}
