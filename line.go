package teq

import "strings"

func lineOf(text string) line {
	return line{text: text}
}

func linesOf(texts ...string) lines {
	result := make(lines, 0, len(texts))
	for _, text := range texts {
		result = append(result, line{text: text})
	}
	return result
}

type line struct {
	indentDepth int
	text        string
}

type lines []line

func (l line) indent() line {
	return line{
		indentDepth: l.indentDepth + 1,
		text:        l.text,
	}
}

func (ls lines) indent() lines {
	result := make(lines, len(ls))
	for i, l := range ls {
		result[i] = l.indent()
	}
	return result
}

func (ls lines) diffSequence() []string {
	result := make([]string, len(ls))
	for i, l := range ls {
		result[i] = strings.Repeat("  ", l.indentDepth) + l.text + "\n"
	}
	return result
}

func (ls lines) ledBy(s string) lines {
	result := ls.clone()
	result[0].text = s + result[0].text
	return result
}

func (ls lines) followedBy(s string) lines {
	result := ls.clone()
	result[len(result)-1].text = result[len(result)-1].text + s
	return result
}

func (ls lines) clone() lines {
	result := make(lines, len(ls))
	copy(result, ls)
	return result
}

func (ls lines) key() string {
	ts := make([]string, len(ls))
	for i, l := range ls {
		ts[i] = l.text
	}
	return strings.Join(ts, "")
}

func keyValue(key lines, value lines) lines {
	result := key.followedBy(": " + value[0].text)
	result = append(result, value[1:]...)
	return result.followedBy(",")
}
