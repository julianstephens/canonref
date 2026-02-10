package bibleref

import (
	"fmt"

	"github.com/julianstephens/canonref/util"
)

// BibleRef represents a reference to a specific passage in the Bible, consisting
// of an OSIS code for the book, a chapter number, and an optional verse or verse range.
type BibleRef struct {
	OSIS    string
	Chapter int
	Verse   *util.VerseRange
}

// String returns a string representation of the BibleRef in the format "OSIS Chapter:Verse"
// or "OSIS Chapter" if Verse is nil.
type Format int

const (
	FormatOSIS      Format = iota // "Prov.31.10-31"
	FormatHuman                   // "Proverbs 31:10–31"
	FormatCanonical               // "Prov 31:10-31"
)

// String returns a string representation in the canonical format,
// e.g. "Prov 3:16" or "Prov 3:16–18" or "Prov 3".
func (r BibleRef) String() string {
	if r.Verse == nil {
		return fmt.Sprintf("%s %d", r.OSIS, r.Chapter)
	}
	return fmt.Sprintf("%s %d:%s", r.OSIS, r.Chapter, r.Verse.String())
}

// Format returns a string representation of the BibleRef in the specified format.
// For FormatOSIS, the format is "OSIS.Chapter.Verse" or "OSIS.Chapter" if Verse is nil.
// For FormatHuman, the format is "BookName Chapter:Verse" or "BookName Chapter" if Verse is nil.
// For FormatCanonical, the format is "OSIS Chapter:Verse" or "OSIS Chapter" if Verse is nil.
func (r BibleRef) Format(f Format, tbl *Table) string {
	switch f {
	case FormatOSIS:
		if r.Verse == nil {
			return fmt.Sprintf("%s.%d", r.OSIS, r.Chapter)
		}
		return fmt.Sprintf("%s.%d.%s", r.OSIS, r.Chapter, r.Verse.String())
	case FormatHuman:
		book := tbl.ByOsis[r.OSIS]
		if r.Verse == nil {
			return fmt.Sprintf("%s %d", book.Name, r.Chapter)
		}
		return fmt.Sprintf("%s %d:%s", book.Name, r.Chapter, r.Verse.String())
	case FormatCanonical:
		if r.Verse == nil {
			return fmt.Sprintf("%s %d", r.OSIS, r.Chapter)
		}
		return fmt.Sprintf("%s %d:%s", r.OSIS, r.Chapter, r.Verse.String())
	default:
		return r.String()
	}
}

// IsChapterOnly returns true if the BibleRef has only a chapter (i.e. it does not have a Verse).
func (r BibleRef) IsChapterOnly() bool {
	return r.Verse == nil
}

// IsSingleVerse returns true if the BibleRef has a single verse
// (i.e. it has a Verse and that Verse does not have an EndVerse).
func (r BibleRef) IsSingleVerse() bool {
	return r.Verse != nil && r.Verse.EndVerse == nil
}

// IsRange returns true if the BibleRef has a verse range
// (i.e. it has a Verse and that Verse has an EndVerse).
func (r BibleRef) IsRange() bool {
	return r.Verse != nil && r.Verse.EndVerse != nil
}

// Validate checks if the BibleRef is valid according to the provided Table.
// It checks if the OSIS code exists in the Table, if the chapter number is valid for the book,
// and if the verse numbers are valid (positive integers and end verse is greater than or equal to start verse).
func (r BibleRef) Validate(tbl *Table) error {
	book, ok := tbl.ByOsis[r.OSIS]
	if !ok {
		return &BibleRefError{
			Kind:    KindUnknownBook,
			Err:     ErrInvalidOSISCode,
			Message: util.Ptr(fmt.Sprintf("unknown OSIS code: %s", r.OSIS)),
		}
	}

	if r.Chapter < 1 || r.Chapter > book.Chapters {
		return &BibleRefError{
			Kind:    KindInvalidChapter,
			Err:     ErrInvalidChapter,
			Message: util.Ptr(fmt.Sprintf("invalid chapter number %d for book %s", r.Chapter, book.Name)),
		}
	}

	if r.Verse != nil {
		if r.Verse.StartVerse < 1 {
			return &BibleRefError{
				Kind:    KindInvalidVerse,
				Err:     ErrInvalidVerse,
				Message: util.Ptr(fmt.Sprintf("start verse must be a positive integer, got %d", r.Verse.StartVerse)),
			}
		}
		if r.Verse.EndVerse != nil {
			if *r.Verse.EndVerse < r.Verse.StartVerse {
				return &BibleRefError{
					Kind:    KindInvalidVerse,
					Err:     ErrInvalidVerse,
					Message: util.Ptr(fmt.Sprintf("end verse must be greater than or equal to start verse, got start: %d, end: %d", r.Verse.StartVerse, *r.Verse.EndVerse)),
				}
			}
		}
	}

	return nil
}

// Book represents a book of the Bible, including its OSIS code,
// name, aliases, testament, order, and number of chapters.
type Book struct {
	OSIS      string
	Name      string
	Aliases   []string
	Testament string
	Order     int
	Chapters  int
}

// Validate checks if the Book has valid data and returns an error if any validation fails.
func (b Book) Validate() error {
	if b.OSIS == "" {
		return &BibleRefError{
			Kind: KindInvalidBook,
			Err:  ErrInvalidOSISCode,
		}
	}

	if b.Name == "" {
		return &BibleRefError{
			Kind:    KindInvalidBook,
			Err:     ErrInvalidBook,
			Message: util.Ptr("book name cannot be empty"),
		}
	}

	if b.Chapters < 1 {
		return &BibleRefError{
			Kind:    KindInvalidBook,
			Err:     ErrInvalidBook,
			Message: util.Ptr("book must have at least one chapter"),
		}
	}

	if b.Order < 1 {
		return &BibleRefError{
			Kind:    KindInvalidBook,
			Err:     ErrInvalidBook,
			Message: util.Ptr("book order must be a positive integer"),
		}
	}

	return nil
}
