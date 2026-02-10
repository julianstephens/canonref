package bibleref

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/julianstephens/canonref/util"
)

// Parse parses a reference string into a BibleRef struct using the provided Table for book lookups.
// It returns a BibleRefError if parsing fails or if the reference is invalid.
func Parse(s string, tbl *Table) (*BibleRef, error) {
	parseResult, err := doParse(s, tbl)
	if err != nil {
		return nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr(fmt.Sprintf("failed to parse reference string: %s", s)),
			Cause:   err,
		}
	}

	return parseResult, nil
}

// MustParse is a helper function that calls Parse and panics if there is an error.
func MustParse(s string, tbl *Table) *BibleRef {
	ref, err := Parse(s, tbl)
	if err != nil {
		panic(fmt.Sprintf("failed to parse reference string: %s, error: %v", s, err))
	}

	return ref
}

func doParse(s string, tbl *Table) (*BibleRef, error) {
	ref, err := parseRefString(s, tbl)
	if err != nil {
		return nil, err
	}

	if err := ref.Validate(tbl); err != nil {
		return nil, err
	}

	return ref, nil
}

func parseRefString(s string, tbl *Table) (*BibleRef, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("reference string cannot be empty"),
		}
	}

	fields := strings.Fields(s)
	if len(fields) < 2 {
		return nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("reference string must contain at least a book and a chapter"),
		}
	}

	bookPart := strings.Join(fields[:len(fields)-1], " ")
	bookStr := NormalizeAlias(bookPart)
	chapterVerseStr, err := parseTail(fields[len(fields)-1])
	if err != nil {
		return nil, err
	}

	bookOsis, ok := tbl.ByAlias[bookStr]
	if !ok {
		bookOsis = bookStr
	}

	book, ok := tbl.ByOsis[bookOsis]
	if !ok {
		return nil, &BibleRefError{
			Kind:    KindUnknownBook,
			Err:     ErrInvalidOSISCode,
			Message: util.Ptr(fmt.Sprintf("unknown book: %s", bookStr)),
		}
	}

	chapter, verseRange, err := parseChapterVerse(chapterVerseStr)
	if err != nil {
		return nil, err
	}
	if verseRange != nil && verseRange.StartVerse < 1 {
		return nil, &BibleRefError{
			Kind:    KindInvalidVerse,
			Err:     ErrInvalidVerse,
			Message: util.Ptr(fmt.Sprintf("invalid verse number: %d", verseRange.StartVerse)),
		}
	}

	ref := &BibleRef{
		OSIS:    book.OSIS,
		Chapter: chapter,
		Verse:   verseRange,
	}
	if err := ref.Validate(tbl); err != nil {
		return nil, err
	}

	return ref, nil
}

func parseChapterVerse(s string) (int, *util.VerseRange, error) {
	parts := strings.Split(s, ":")
	if len(parts) == 0 {
		return 0, nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("chapter and verse string must contain at least a chapter"),
		}
	}
	if len(parts) > 2 {
		return 0, nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("chapter and verse string must contain at most one colon"),
		}
	}

	chapterStr := parts[0]
	chapter, err := strconv.Atoi(chapterStr)
	if err != nil {
		return 0, nil, &BibleRefError{
			Kind:    KindInvalidChapter,
			Err:     ErrInvalidChapter,
			Message: util.Ptr(fmt.Sprintf("invalid chapter: %s", chapterStr)),
			Cause:   err,
		}
	}

	if len(parts) == 1 {
		return chapter, nil, nil
	}

	verseStr := NormalizeVerseRange(parts[1])

	if strings.Contains(verseStr, util.EnDash) {
		verseParts := strings.Split(verseStr, util.EnDash)
		verseRange, err := parseVerseRange(verseStr, verseParts)
		if err != nil {
			return 0, nil, err
		}
		return chapter, verseRange, nil
	} else {
		startVerse, err := strconv.Atoi(verseStr)
		if err != nil {
			return 0, nil, &BibleRefError{
				Kind:    KindInvalidVerse,
				Err:     ErrInvalidVerse,
				Message: util.Ptr(fmt.Sprintf("invalid verse: %s", verseStr)),
				Cause:   err,
			}
		}
		return chapter, &util.VerseRange{StartVerse: startVerse}, nil
	}
}

func parseVerseRange(s string, parts []string) (*util.VerseRange, error) {
	if len(parts) != 2 {
		return nil, &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr(fmt.Sprintf("invalid verse range: %s", s)),
		}
	}

	startVerse, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, &BibleRefError{
			Kind:    KindInvalidVerse,
			Err:     ErrInvalidVerse,
			Message: util.Ptr(fmt.Sprintf("invalid start verse: %s", parts[0])),
			Cause:   err,
		}
	}

	endVerse, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, &BibleRefError{
			Kind:    KindInvalidVerse,
			Err:     ErrInvalidVerse,
			Message: util.Ptr(fmt.Sprintf("invalid end verse: %s", parts[1])),
			Cause:   err,
		}
	}

	return &util.VerseRange{StartVerse: startVerse, EndVerse: &endVerse}, nil
}

func parseTail(tail string) (string, error) {
	if tail == "" {
		return "", &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr("tail cannot be empty"),
		}
	}

	if tail[0] < '0' || tail[0] > '9' {
		return "", &BibleRefError{
			Kind:    KindInvalidChapter,
			Err:     ErrInvalidChapter,
			Message: util.Ptr(fmt.Sprintf("tail must start with a digit, got: %s", tail)),
		}
	}

	i := 0
	for i < len(tail) && tail[i] >= '0' && tail[i] <= '9' {
		i++
	}

	if i == len(tail) {
		return tail, nil
	}

	if tail[i] != ':' {
		return "", &BibleRefError{
			Kind:    KindParse,
			Err:     ErrBibleRefParseFailed,
			Message: util.Ptr(fmt.Sprintf("expected ':' after chapter, got: %c at position %d", tail[i], i)),
		}
	}

	versesPart := tail[i+1:]
	if versesPart == "" {
		return "", &BibleRefError{
			Kind:    KindInvalidVerse,
			Err:     ErrInvalidVerse,
			Message: util.Ptr("verse part cannot be empty after ':'"),
		}
	}

	normalizedVerses := NormalizeVerseRange(versesPart)

	return tail[:i] + ":" + normalizedVerses, nil
}

// NormalizeAlias normalizes a book name or alias by trimming whitespace, converting to lowercase,
// removing punctuation, and replacing hyphens with en dashes. It also handles common roman numeral prefixes.
func NormalizeAlias(s string) string {
	res := strings.TrimSpace(s)
	res = strings.ToLower(res)
	res = strings.ReplaceAll(res, ".", "")
	res = strings.ReplaceAll(res, util.EnDash, util.Hyphen)

	// handle roman numeral prefixes
	res = strings.ReplaceAll(res, "iii ", "3 ")
	res = strings.ReplaceAll(res, "ii ", "2 ")
	res = strings.ReplaceAll(res, "i ", "1 ")

	// unicode apostrophes & quotation marks
	res = strings.ReplaceAll(res, "’", "'")
	res = strings.ReplaceAll(res, "‘", "'")
	res = strings.ReplaceAll(res, "“", "\"")
	res = strings.ReplaceAll(res, "”", "\"")

	return res
}

// NormalizeVerseRange normalizes a verse range string by trimming whitespace,
// replacing hyphens with en dashes, and removing spaces.
func NormalizeVerseRange(s string) string {
	res := strings.TrimSpace(s)
	res = strings.ReplaceAll(res, util.Hyphen, util.EnDash)
	res = strings.ReplaceAll(res, " ", "")
	return res
}
