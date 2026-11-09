package search_engine

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type IndexStats interface {
	DF(term string) int
}

func Compile(n Node, ix *Index) (Iterator, error) {
	opt, err := Optimize(n, ix)
	if err != nil {
		return nil, err
	}
	return build(opt, ix), nil
}

func Optimize(n Node, st IndexStats) (Node, error) {
	out, _, err := optimize(n, st)
	return out, err
}

type scored struct {
	node Node
	cost int
}

func optimize(n Node, st IndexStats) (Node, int, error) {
	switch n := n.(type) {
	case nil:
		return nil, 0, errors.New("узел запроса отсутствует (nil)")

	case TermNode:
		t, err := normalizeTerm(n.Text)
		if err != nil {
			return nil, 0, err
		}
		return Term(t), st.DF(t), nil

	case AndNode:
		kids, err := optimizeKids("And", n.Kids, st)
		if err != nil {
			return nil, 0, err
		}
		sort.SliceStable(kids, func(i, j int) bool { return kids[i].cost < kids[j].cost })
		return And(nodes(kids)...), kids[0].cost, nil

	case OrNode:
		kids, err := optimizeKids("Or", n.Kids, st)
		if err != nil {
			return nil, 0, err
		}
		total := 0
		for _, k := range kids {
			total += k.cost
		}
		return Or(nodes(kids)...), total, nil

	case AndNotNode:
		incl, cost, err := optimize(n.Include, st)
		if err != nil {
			return nil, 0, fmt.Errorf("в include узла AndNot: %w", err)
		}
		excl, _, err := optimize(n.Exclude, st)
		if err != nil {
			return nil, 0, fmt.Errorf("в exclude узла AndNot: %w", err)
		}
		return AndNot(incl, excl), cost, nil

	default:
		return nil, 0, fmt.Errorf("неизвестный тип узла %T", n)
	}
}

func optimizeKids(op string, kids []Node, st IndexStats) ([]scored, error) {
	if len(kids) == 0 {
		return nil, fmt.Errorf("%s без дочерних узлов", op)
	}
	out := make([]scored, len(kids))
	for i, k := range kids {
		node, cost, err := optimize(k, st)
		if err != nil {
			return nil, fmt.Errorf("в узле %s (дочерний №%d): %w", op, i+1, err)
		}
		out[i] = scored{node, cost}
	}
	return out, nil
}

func nodes(s []scored) []Node {
	out := make([]Node, len(s))
	for i, k := range s {
		out[i] = k.node
	}
	return out
}

func normalizeTerm(text string) (string, error) {
	if strings.ContainsFunc(strings.TrimSpace(text), unicode.IsSpace) {
		return "", fmt.Errorf("терм %q содержит пробел: слова хранятся в индексе по отдельности, "+
			"запрос из нескольких слов нужно строить через And/Or", text)
	}
	t := QueryTerm(text)
	if t == "" {
		return "", fmt.Errorf("терм %q пуст после нормализации (в нём нет букв и цифр)", text)
	}
	return t, nil
}

func build(n Node, ix *Index) Iterator {
	switch n := n.(type) {
	case TermNode:
		return ix.Postings(n.Text)
	case AndNode:
		return IterAnd(buildAll(n.Kids, ix)...)
	case OrNode:
		return IterOr(buildAll(n.Kids, ix)...)
	case AndNotNode:
		return IterAndNot(build(n.Include, ix), build(n.Exclude, ix))
	}
	return Empty()
}

func buildAll(kids []Node, ix *Index) []Iterator {
	out := make([]Iterator, len(kids))
	for i, k := range kids {
		out[i] = build(k, ix)
	}
	return out
}
