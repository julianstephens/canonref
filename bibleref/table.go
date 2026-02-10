package bibleref

import (
	"encoding/json"

	"github.com/julianstephens/canonref/util"
)

// Table represents a mapping of OSIS codes to Books and aliases to OSIS codes.
type Table struct {
	ByOsis  map[string]Book
	ByAlias map[string]string
}

// NewTable creates a new Table from a slice of Books.
// It validates each Book and returns an error if any Book is invalid.
func NewTable(books []Book) (*Table, error) {
	tbl := &Table{
		ByOsis:  make(map[string]Book, len(books)),
		ByAlias: make(map[string]string, len(books)),
	}

	for _, book := range books {
		if err := book.Validate(); err != nil {
			return nil, err
		}
		tbl.ByOsis[book.OSIS] = book
		for _, alias := range book.Aliases {
			tbl.ByAlias[alias] = book.OSIS
		}
	}

	return tbl, nil
}

// LoadTableFromJSON loads a Table from JSON data.
// The JSON should be an array of Book objects.
func LoadTableFromJSON(jsonData []byte) (*Table, error) {
	var books []Book
	if err := json.Unmarshal(jsonData, &books); err != nil {
		return nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("failed to parse JSON data"),
			Cause:   err,
		}
	}

	return NewTable(books)
}
