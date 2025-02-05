package doc

func Block(open inline, contents []Doc, close inline) *block {
	return &block{open, contents, close}
}

type block struct {
	open     Doc
	contents []Doc
	close    Doc
}

func (b block) print(depth int, buf *buffer) {
	b.open.print(depth, buf)
	for _, d := range b.contents {
		d.print(depth+1, buf)
	}
	b.close.print(depth, buf)
}

func (b block) Left() Doc {
	var c block
	c.open = b.open.Left()
	for _, d := range b.contents {
		c.contents = append(c.contents, d.Left())
	}
	c.close = c.close.Left()
	return c
}

func (b block) Right() Doc {
	var c block
	c.open = c.open.Right()
	for _, d := range b.contents {
		c.contents = append(c.contents, d.Right())
	}
	c.close = c.close.Right()
	return c
}

func (b block) AddPrefix(prefix string) Doc {
	b.open = b.open.AddPrefix(prefix)
	return b
}

func (b block) AddSuffix(suffix string) Doc {
	b.close = b.close.AddSuffix(suffix)
	return b
}
