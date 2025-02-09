package doc

import (
	"strings"
)

const (
	indentSize = 2
	sightWidth = 2
)

type Doc interface {
	print(depth int) virtualLines
	Left() Doc
	Right() Doc
	AddPrefix(prefix string) Doc
	AddSuffix(suffix string) Doc
}

var _ Doc = block{}
var _ Doc = inline{}

func PrintDoc(ds []Doc) string {
	vls := virtualLines{}
	for _, d := range ds {
		vls = append(vls, d.print(0)...)
	}

	shouldPrint := make([]bool, len(vls))
	for i, vl := range vls {
		if vl.isContext {
			shouldPrint[i] = true
		}
		if vl.isDiff {
			for j := i - sightWidth; j <= i+sightWidth; j++ {
				if j >= 0 && j < len(vls) {
					shouldPrint[j] = true
				}
			}
		}
	}

	var lines []string
	for i, vl := range vls {
		if !shouldPrint[i] {
			continue
		}
		if i > 0 && !shouldPrint[i-1] {
			lines = append(lines, ":")
		}
		lines = append(lines, vl.text)
	}

	return strings.Join(lines, "\n")
}

type virtualLine struct {
	isContext bool
	isDiff    bool
	text      string
}

type virtualLines []virtualLine

func (v virtualLines) asContext() virtualLines {
	var out virtualLines
	for _, l := range v {
		l.isContext = true
		out = append(out, l)
	}
	return out
}

func (v virtualLines) hasDiff() bool {
	for _, l := range v {
		if l.isDiff {
			return true
		}
	}
	return false
}
