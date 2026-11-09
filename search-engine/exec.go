package search_engine

func Drain(it Iterator) []uint32 {
	var out []uint32
	for d := it.CurDocID(); d != NoMore; d = it.Next() {
		out = append(out, d)
	}
	return out
}

// EmptyIterator -- returns nothing, using at empty OR()/AND()

type EmptyIter struct{}

func (EmptyIter) CurDocID() uint32 {
	return NoMore
}
func (EmptyIter) Next() uint32 {
	return NoMore
}
func (EmptyIter) Advance(uint32) uint32 {
	return NoMore
}

func Empty() Iterator {
	return EmptyIter{}
}

// Boolean AND iterator, private func process find pos where all its agree at doc

type IterAND struct {
	its []Iterator
	cur uint32
}

func IterAnd(its ...Iterator) Iterator {
	switch len(its) {
	case 0:
		return Empty()
	case 1:
		return its[0]
	}

	andRes := &IterAND{its: its}
	andRes.cur = andRes.process(0)
	return andRes
}

func (a *IterAND) process(v uint32) uint32 {
	for {
		v = a.its[0].Advance(v)
		if v == NoMore {
			return NoMore
		}
		agreed := true
		for _, it := range a.its[1:] {
			d := it.Advance(v)
			if d == NoMore {
				return NoMore
			}
			if d > v {
				v, agreed = d, false
				break
			}
		}
		if agreed {
			return v
		}
	}
}

func (a *IterAND) CurDocID() uint32 {
	return a.cur
}

func (a *IterAND) Next() uint32 {
	if a.cur == NoMore {
		return NoMore
	}
	a.cur = a.process(a.cur + 1)
	return a.cur
}

func (a *IterAND) Advance(v uint32) uint32 {
	if a.cur >= v { // включает случай NoMore
		return a.cur
	}
	a.cur = a.process(v)
	return a.cur
}

// Boolean OR iterator, private func process find the smallest doc at its

type IterOR struct {
	its []Iterator
	cur uint32
}

func IterOr(its ...Iterator) Iterator {
	switch len(its) {
	case 0:
		return Empty()
	case 1:
		return its[0]
	}

	orRes := &IterOR{its: its}
	orRes.process()
	return orRes
}

func (o *IterOR) process() {
	m := NoMore
	for _, it := range o.its {
		if d := it.CurDocID(); d < m {
			m = d
		}
	}
	o.cur = m
}

func (o *IterOR) CurDocID() uint32 { return o.cur }

func (o *IterOR) Next() uint32 {
	if o.cur == NoMore {
		return NoMore
	}
	for _, it := range o.its {
		if it.CurDocID() == o.cur {
			it.Next()
		}
	}
	o.process()
	return o.cur
}

func (o *IterOR) Advance(v uint32) uint32 {
	if o.cur >= v {
		return o.cur
	}
	for _, it := range o.its {
		it.Advance(v)
	}
	o.process()
	return o.cur
}

// Boolean AND NOT iterator, func private process find the first acceptable

type IterANDNOT struct {
	incl Iterator
	excl Iterator
	cur  uint32
}

func IterAndNot(incl, excl Iterator) Iterator {
	n := &IterANDNOT{incl: incl, excl: excl}
	n.process()
	return n
}

func (n *IterANDNOT) process() {
	for {
		d := n.incl.CurDocID()
		if d == NoMore {
			n.cur = NoMore
			return
		}
		if n.excl.Advance(d) != d {
			n.cur = d
			return
		}
		n.incl.Next()
	}
}

func (n *IterANDNOT) CurDocID() uint32 { return n.cur }

func (n *IterANDNOT) Next() uint32 {
	if n.cur == NoMore {
		return NoMore
	}
	n.incl.Next()
	n.process()
	return n.cur
}

func (n *IterANDNOT) Advance(v uint32) uint32 {
	if n.cur >= v {
		return n.cur
	}
	n.incl.Advance(v)
	n.process()
	return n.cur
}
