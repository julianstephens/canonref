package bibleref

import (
	"encoding/json"

	"github.com/julianstephens/canonref/util"
)

// booksWrapper is used to unmarshal JSON with schema and work fields
type booksWrapper struct {
	Schema int    `json:"schema"`
	Work   string `json:"work"`
	Books  []Book `json:"books"`
}

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
			normalizedAlias := NormalizeAlias(alias)
			tbl.ByAlias[normalizedAlias] = book.OSIS
		}
		if !contains(tbl.ByAlias, NormalizeAlias(book.OSIS)) {
			tbl.ByAlias[NormalizeAlias(book.OSIS)] = book.OSIS
		}
	}

	return tbl, nil
}

// LoadTableFromJSON loads a Table from JSON data.
// The JSON should have schema, work, and books fields with an array of Book objects.
func LoadTableFromJSON(jsonData []byte) (*Table, error) {
	var wrapper booksWrapper
	if err := json.Unmarshal(jsonData, &wrapper); err != nil {
		return nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("failed to parse JSON data"),
			Cause:   err,
		}
	}

	return NewTable(wrapper.Books)
}

func contains(m map[string]string, key string) bool {
	_, exists := m[key]
	return exists
}
