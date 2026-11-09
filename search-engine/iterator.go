package search_engine

import (
	"fmt"
	"math"
	"sort"
)

const NoMore uint32 = math.MaxUint32

type ListIterator struct {
	list []uint32
	pos  int
}

func NewListIterator(elems ...uint32) *ListIterator {
	return &ListIterator{
		list: elems,
	}
}

type Iterator interface {
	CurDocID() uint32
	Next() uint32
	Advance(v uint32) uint32
}

func (it *ListIterator) String() string {
	if len(it.list) > 10 {
		return fmt.Sprintf("ListIterator pos:%d; list:%v", it.pos, it.list[:10])
	}
	return fmt.Sprintf("ListIterator pos:%d; list:%v", it.pos, it.list)
}

func (it *ListIterator) CurDocID() uint32 {
	if it.pos >= len(it.list) {
		return NoMore
	}
	return it.list[it.pos]
}

func (it *ListIterator) Next() uint32 {
	if it.pos < len(it.list) {
		it.pos++
	}
	return it.CurDocID()
}

func (it *ListIterator) Advance(v uint32) uint32 {
	if it.pos >= len(it.list) {
		return NoMore
	}
	if it.list[it.pos] >= v {
		return it.list[it.pos]
	}
	rest := it.list[it.pos+1:]
	k := sort.Search(len(rest), func(i int) bool { return rest[i] >= v })
	it.pos += 1 + k
	return it.CurDocID()
}
