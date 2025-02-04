package diff

import "strings"

const indentSize = 2

type line struct {
	onLeft  bool
	onRight bool
	text    string
	isOpen  bool
	isClose bool
}

func leftLine(text string) line {
	return line{onLeft: true, text: text}
}
func rightLine(text string) line {
	return line{onRight: true, text: text}
}
func bothLine(text string) line {
	return line{onLeft: true, onRight: true, text: text}
}

func (l line) open() line {
	l.isOpen = true
	return l
}
func (l line) close() line {
	l.isClose = true
	return l
}

func (l line) overrideText(text string) line {
	l.text = text
	return l
}

func (l line) print(depth int) string {
	var marker string = "  "
	if l.onLeft && !l.onRight {
		marker = "- "
	} else if !l.onLeft && l.onRight {
		marker = "+ "
	}
	return marker + strings.Repeat(" ", depth*indentSize) + l.text
}

type lines []line

func (l *lines) add(line line) {
	*l = append(*l, line)
}

func (l *lines) concat(other lines) {
	*l = append(*l, other...)
}

func (l *lines) left() {
	for i := range *l {
		(*l)[i].onLeft = true
		(*l)[i].onRight = false
	}
}

func (l *lines) right() {
	for i := range *l {
		(*l)[i].onLeft = false
		(*l)[i].onRight = true
	}
}

func (l lines) print() string {
	var depth int = 0
	var ls []string
	for _, ln := range l {
		if ln.isClose {
			depth--
		}
		ls = append(ls, ln.print(depth))
		if ln.isOpen {
			depth++
		}
	}
	return strings.Join(ls, "\n")
}
