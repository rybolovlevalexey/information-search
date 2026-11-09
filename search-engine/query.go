package search_engine

import "strings"

type Node interface {
	String() string
	node()
}

type TermNode struct {
	Text string
}

type AndNode struct {
	Kids []Node
}

type OrNode struct {
	Kids []Node
}

type AndNotNode struct {
	Include Node
	Exclude Node
}

func (TermNode) node()   {}
func (AndNode) node()    {}
func (OrNode) node()     {}
func (AndNotNode) node() {}

func Term(text string) Node {
	return TermNode{Text: text}
}

func And(kids ...Node) Node {
	return AndNode{Kids: append([]Node(nil), kids...)}
}

func Or(kids ...Node) Node {
	return OrNode{Kids: append([]Node(nil), kids...)}
}

func AndNot(include, exclude Node) Node {
	return AndNotNode{Include: include, Exclude: exclude}
}

func str(n Node) string {
	if n == nil {
		return "<nil>"
	}
	return n.String()
}

func (t TermNode) String() string {
	return t.Text
}

func (a AndNode) String() string {
	return join(a.Kids, " AND ")
}

func (o OrNode) String() string {
	return join(o.Kids, " OR ")
}

func (n AndNotNode) String() string {
	return "(" + str(n.Include) + " ANDNOT " + str(n.Exclude) + ")"
}

func join(kids []Node, sep string) string {
	parts := make([]string, len(kids))
	for i, k := range kids {
		parts[i] = str(k)
	}
	return "(" + strings.Join(parts, sep) + ")"
}
