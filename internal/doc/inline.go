package doc

import "strings"

func LeftInline(text string) inline {
	return inline{onLeft: true, text: text}
}

func RightInline(text string) inline {
	return inline{onRight: true, text: text}
}

func BothInline(text string) inline {
	return inline{onRight: true, onLeft: true, text: text}
}

type inline struct {
	onLeft  bool
	onRight bool
	text    string
}

func (l inline) print(depth int, buf *buffer) {
	var marker string = "  "
	if l.onLeft && !l.onRight {
		marker = "- "
	} else if !l.onLeft && l.onRight {
		marker = "+ "
	}
	buf.push(
		marker + strings.Repeat(" ", depth*indentSize) + l.text,
	)
}

func (l inline) Left() Doc {
	l.onLeft = true
	l.onRight = false
	return l
}

func (l inline) Right() Doc {
	l.onLeft = false
	l.onRight = true
	return l
}

func (l inline) AddPrefix(prefix string) Doc {
	l.text = prefix + l.text
	return l
}

func (l inline) AddSuffix(suffix string) Doc {
	l.text = l.text + suffix
	return l
}
