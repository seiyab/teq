package diff

import "strings"

const indentSize = 2

type line struct {
	marker string
	depth  int
	text   string
}

func leftLine(text string) line {
	return line{marker: "- ", text: text}
}
func rightLine(text string) line {
	return line{marker: "+ ", text: text}
}
func bothLine(text string) line {
	return line{text: text}
}

func (l line) indent() line {
	l.depth++
	return l
}

func (l line) print() string {
	return l.marker + strings.Repeat(" ", l.depth*indentSize) + l.text
}

type lines []line

func (l *lines) add(line line) {
	*l = append(*l, line)
}

func (l *lines) concat(other lines) {
	*l = append(*l, other...)
}

func (l lines) indent() lines {
	result := lines{}
	for _, line := range l {
		result.add(line.indent())
	}
	return result
}

func (l lines) print() string {
	ls := make([]string, 0, len(l))
	for _, line := range l {
		ls = append(ls, line.print())
	}
	return strings.Join(ls, "\n")
}
