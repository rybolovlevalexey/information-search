package search_engine

import (
	"sort"
)

type ScoreFunc func(pos int, d *Doc) float64

func l0Score(pos int, _ *Doc) float64 {
	return -1 * float64(pos)
}

type Index struct {
	docs     []Doc
	postings map[string][]uint32
}

type InfoStats struct {
	Docs     int
	Terms    int
	Postings int
}

func Build(docs []Doc) *Index {
	return BuildWith(docs, l0Score)
}

func BuildWith(docs []Doc, score ScoreFunc) *Index {
	ordered := orderByScore(docs, score)
	ix := &Index{docs: ordered, postings: make(map[string][]uint32)}
	for i := range ordered {
		id := uint32(i)
		elem := ordered[i]
		for _, t := range Tokenize(elem.Text) {
			l := ix.postings[t]
			if n := len(l); n > 0 && l[n-1] == id {
				continue
			}
			ix.postings[t] = append(l, id)
		}
	}
	return ix
}

func orderByScore(docs []Doc, score ScoreFunc) []Doc {
	type scored struct {
		pos int
		s   float64
	}
	sc := make([]scored, len(docs))
	for i := range docs {
		sc[i] = scored{i, score(i, &docs[i])}
	}
	sort.SliceStable(sc, func(a, b int) bool { return sc[a].s > sc[b].s })
	out := make([]Doc, len(docs))
	for i, s := range sc {
		out[i] = docs[s.pos]
	}
	return out
}

func (ix *Index) NumDocs() int { return len(ix.docs) }

func (ix *Index) Doc(id uint32) *Doc {
	if int(id) >= len(ix.docs) {
		return nil
	}
	return &ix.docs[id]
}

func (ix *Index) Postings(term string) *ListIterator {
	return &ListIterator{list: ix.postings[term]}
}

func (ix *Index) GetStatics() InfoStats {
	s := InfoStats{Docs: len(ix.docs), Terms: len(ix.postings)}
	for _, l := range ix.postings {
		s.Postings += len(l)
	}
	return s
}

func (ix *Index) DF(term string) int {
	res, ok := ix.postings[term]
	if !ok {
		return 0
	}
	return len(res)
}
