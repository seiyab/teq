package doc

func Block(open inline, contents []Doc, close inline) *block {
	return &block{open, contents, close}
}

type block struct {
	open     Doc
	contents []Doc
	close    Doc
}

func (b block) print(depth int) virtualLines {
	o := b.open.print(depth)
	c := b.close.print(depth)

	var body virtualLines
	for _, d := range b.contents {
		ls := d.print(depth + 1)
		body = append(body, ls...)
	}

	if body.hasDiff() {
		o = o.asContext()
		c = c.asContext()
	}
	return append(append(o, body...), c...)
}

func (b block) Left() Doc {
	var cs []Doc
	cs, b.contents = b.contents, nil
	b.open = b.open.Left()
	for _, d := range cs {
		b.contents = append(b.contents, d.Left())
	}
	b.close = b.close.Left()
	return b
}

func (b block) Right() Doc {
	var cs []Doc
	cs, b.contents = b.contents, nil
	b.open = b.open.Right()
	for _, d := range cs {
		b.contents = append(b.contents, d.Right())
	}
	b.close = b.close.Right()
	return b
}

func (b block) AddPrefix(prefix string) Doc {
	b.open = b.open.AddPrefix(prefix)
	return b
}

func (b block) AddSuffix(suffix string) Doc {
	b.close = b.close.AddSuffix(suffix)
	return b
}
