package search_engine

// Правила очистки и разбиения на термы
// - слова берутся по пробелам + очистка от пунктуации
// - точка между цифрами (например, номер спецификации 5.3.5.1) должна давать один терм
// - буквы с цифрами не разделяются 5G, SIB1, v1760
// - всё приводится к нижнему регистру

import (
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	var out []string
	for _, w := range strings.Fields(text) {
		out = appendWord(out, w)
	}
	return out
}

func QueryTerm(query string) string {
	query = trimWord(query)
	queryList := make([]string, 0)
	for _, elem := range strings.Fields(query) {
		queryList = append(queryList, strings.ToLower(elem))
	}
	return strings.Join(queryList, " ")
}

func trimWord(w string) string {
	w = strings.TrimFunc(w, notAlnum)
	if strings.HasSuffix(w, "'s") {
		w = strings.TrimFunc(strings.TrimSuffix(w, "'s"), notAlnum)
	}
	return w
}

func appendWord(out []string, word string) []string {
	word = trimWord(word)
	if word == "" {
		return out
	}
	start := len(out)

	addToList := func(t string) {
		t = strings.ToLower(t)
		if t == "" {
			return
		}
		for _, elem := range out[start:] {
			if elem == t {
				return
			}
		}
		out = append(out, t)
	}

	if validateWord(word) {
		addToList(word)
	}

	return out
}

func isAlnum(r rune) bool  { return unicode.IsLetter(r) || unicode.IsDigit(r) }
func notAlnum(r rune) bool { return !isAlnum(r) }

func validateWord(w string) bool {
	for _, r := range w {
		if !isAlnum(r) {
			return false
		}
	}
	return true
}
