package search_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Doc struct {
	ID        string   `json:"id"`
	Spec      string   `json:"spec"`
	SpecTitle string   `json:"spec_title"`
	Version   string   `json:"version"`
	Clause    string   `json:"clause"`
	Path      []string `json:"path"`
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	Part      int      `json:"part"`
}

func Load(path string) ([]Doc, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

func Read(r io.Reader) ([]Doc, error) {
	dec := json.NewDecoder(r)
	var docs []Doc
	for {
		var d Doc
		err := dec.Decode(&d)
		if errors.Is(err, io.EOF) {
			return docs, nil
		}
		if err != nil {
			return nil, fmt.Errorf("запись %d: %w", len(docs)+1, err)
		}
		docs = append(docs, d)
	}
}
