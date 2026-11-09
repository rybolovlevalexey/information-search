// go run main.go -corpus ../../data/prep-specs-1.jsonl
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	se "search-engine"
)

const snippetRunes = 200

type demoQuery struct {
	Title string
	Query se.Node
}

func main() {
	corpusPath := flag.String("corpus", "../data/corpus.jsonl", "путь к corpus.jsonl")
	limit := flag.Int("limit", 5, "сколько результатов показывать для каждого запроса")
	flag.Parse()

	if err := run(os.Stdout, *corpusPath, *limit, demoQueries); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

func run(w io.Writer, corpusPath string, limit int, queries []demoQuery) error {
	t0 := time.Now()
	docs, err := se.Load(corpusPath)
	if err != nil {
		return fmt.Errorf("загрузка корпуса: %w", err)
	}
	tLoad := time.Since(t0)

	t0 = time.Now()
	ix := se.Build(docs)
	tBuild := time.Since(t0)

	st := ix.GetStatics()
	fmt.Fprintf(w, "корпус:   %d документов, загрузка %v\n", st.Docs, tLoad.Round(time.Millisecond))
	fmt.Fprintf(w, "индекс:   %d термов, %d записей в posting lists, построение %v\n\n",
		st.Terms, st.Postings, tBuild.Round(time.Millisecond))

	for i, dq := range queries {
		fmt.Fprintf(w, "=== %d. %s\n", i+1, dq.Title)
		fmt.Fprintf(w, "запрос:  %s\n", dq.Query)

		plan, err := se.Optimize(dq.Query, ix)
		if err != nil {
			fmt.Fprintf(w, "ошибка:  %v\n\n", err)
			continue
		}
		fmt.Fprintf(w, "план:    %s\n", plan)

		t0 := time.Now()
		it, err := se.Compile(dq.Query, ix)
		if err != nil {
			fmt.Fprintf(w, "ошибка:  %v\n\n", err)
			continue
		}
		ids := se.Drain(it)
		elapsed := time.Since(t0)

		shown := len(ids)
		if shown > limit {
			shown = limit
		}
		fmt.Fprintf(w, "найдено: %d (показано %d), время %v\n", len(ids), shown, elapsed.Round(time.Microsecond))
		for n, id := range ids[:shown] {
			d := ix.Doc(id)
			fmt.Fprintf(w, "%3d. [%s %s] %s\n", n+1, d.Spec, d.Clause, strings.Join(d.Path, " > "))
			fmt.Fprintf(w, "     %s\n", snippet(d.Text, snippetRunes))
		}
		fmt.Fprintln(w)
	}
	return nil
}

func snippet(text string, maxRunes int) string {
	s := strings.Join(strings.Fields(text), " ")
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	return string([]rune(s)[:maxRunes]) + "…"
}

var demoQueries = []demoQuery{
	//{
	//	Title: "3gpp and not 5g",
	//	Query: se.AndNot(se.Term("3GPP"), se.Term("5G")),
	//},
	//{
	//	Title: "4g and not 5g",
	//	Query: se.AndNot(se.Term("4g"), se.Term("5g")),
	//},

	{
		Title: "",
		Query: se.AndNot(se.Term("And"), se.Term("And")),
	},

	//{
	//	Title: "4g data - 1",
	//	Query: se.And(se.Term("4g"), se.Term("MME")),
	//},
	//{
	//	Title: "4g data - 2",
	//	Query: se.Or(se.Term("4g"), se.Term("MME")),
	//},
	//{
	//	Title: "4g data - 3",
	//	Query: se.AndNot(se.Term("4g"), se.Term("MME")),
	//},

	//{
	//	Title: "Один терм",
	//	Query: se.Term("RRCReconfiguration"),
	//},
	//{
	//	Title: "Пересечение двух термов",
	//	Query: se.And(se.Term("handover"), se.Term("conditional")),
	//},
	//{
	//	Title: "Объединение трёх термов",
	//	Query: se.Or(se.Term("SRB1"), se.Term("SRB2"), se.Term("SRB3")),
	//},
	//{
	//	Title: "AndNot",
	//	Query: se.AndNot(se.And(se.Term("RRC"), se.Term("reconfiguration")), se.Term("SCG")),
	//},
	//{
	//	Title: "Комбинация And и Or",
	//	Query: se.And(se.Term("UE"), se.Or(se.Term("T304"), se.Term("T311")), se.Term("expiry")),
	//},
}
