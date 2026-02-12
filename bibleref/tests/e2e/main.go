package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/julianstephens/canonref/bibleref"
	"github.com/julianstephens/canonref/util"
)

type Suite struct {
	tbl        *bibleref.Table
	testInputs []struct {
		input    string
		expected string
	}
}

func NewSuite(inputs []struct {
	input    string
	expected string
}, bookPath string) (*Suite, error) {
	data, err := os.ReadFile(bookPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read books.json: %v", err)
	}

	tbl, err := bibleref.LoadTableFromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load table from JSON: %v", err)
	}

	return &Suite{tbl: tbl, testInputs: inputs}, nil
}

func (s *Suite) TestParse() error {
	for _, test := range s.testInputs {
		if err := s.runParseTest(test.input, test.expected); err != nil {
			return fmt.Errorf("test failed for input '%s': %v", test.input, err)
		}
	}
	return nil
}

func (s *Suite) runParseTest(input, expected string) error {
	ref, err := bibleref.Parse(input, s.tbl)
	if err != nil {
		return err
	}

	if ref.Format(bibleref.FormatCanonical, nil) != expected {
		return fmt.Errorf("expected '%s', got '%s'", expected, ref.Format(bibleref.FormatCanonical, nil))
	}

	return nil
}

func main() {
	inputs := []struct {
		input    string
		expected string
	}{
		{"Proverbs 31:10-31", fmt.Sprintf("Prov 31:10%s31", util.EnDash)},
		{"Genesis 1:1", "Gen 1:1"},
		{"II Kings 20", "2 Kgs 20"},
		{"Ps 119:105", "Ps 119:105"},
		{"lam 1:1", "Lam 1:1"},
		{"The WISDOM of soloMon 2", "Wis 2"},
		{"  Col   4 ", "Col 4"},
	}

	bookPath := flag.String("bookPath", "./books.json", "The path to the generated books.json to build the table from")
	flag.Parse()

	s, err := NewSuite(inputs, *bookPath)
	if err != nil {
		println("Failed to set up test suite:", err.Error())
		return
	}

	if err := s.TestParse(); err != nil {
		println("Test failed during parsing:", err.Error())
		return
	}

	println("All tests passed!")
}
