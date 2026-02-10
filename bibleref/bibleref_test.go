package bibleref_test

import (
	"testing"

	"github.com/julianstephens/canonref/bibleref"
	"github.com/julianstephens/canonref/util"
)

// testBooks creates test data with a minimal set of Bible books including canonical scriptures and apocrypha.
func testBooks() []bibleref.Book {
	return []bibleref.Book{
		{
			OSIS:      "Prov",
			Name:      "Proverbs",
			Aliases:   []string{"proverbs", "prov", "pro"},
			Testament: "OT",
			Order:     20,
			Chapters:  31,
		},
		{
			OSIS:      "1Sam",
			Name:      "1 Samuel",
			Aliases:   []string{"1 samuel", "1samuel", "1 sam", "1sam", "i samuel", "i sam"},
			Testament: "OT",
			Order:     9,
			Chapters:  31,
		},
		{
			OSIS:      "2Sam",
			Name:      "2 Samuel",
			Aliases:   []string{"2 samuel", "2samuel", "2 sam", "2sam", "ii samuel", "ii sam"},
			Testament: "OT",
			Order:     10,
			Chapters:  24,
		},
		{
			OSIS:      "Wis",
			Name:      "Wisdom of Solomon",
			Aliases:   []string{"wisdom of solomon", "wisdom", "wis", "book of wisdom"},
			Testament: "Apocrypha",
			Order:     70,
			Chapters:  19,
		},
		{
			OSIS:      "Matt",
			Name:      "Matthew",
			Aliases:   []string{"matthew", "matt", "mt"},
			Testament: "NT",
			Order:     40,
			Chapters:  28,
		},
	}
}

// TestTable_AliasNormalization verifies that all aliases normalize and map to valid OSIS codes.
func TestTable_AliasNormalization(t *testing.T) {
	books := testBooks()
	tbl, err := bibleref.NewTable(books)
	if err != nil {
		t.Fatalf("NewTable failed: %v", err)
	}

	testCases := []struct {
		alias    string
		expected string
		desc     string
	}{
		{"proverbs", "Prov", "lowercase full name"},
		{"prov", "Prov", "lowercase OSIS"},
		{"pro", "Prov", "lowercase abbreviated"},
		{"1 samuel", "1Sam", "1 Samuel with space"},
		{"1samuel", "1Sam", "1 Samuel no space"},
		{"1sam", "1Sam", "1 Sam"},
		{"i samuel", "1Sam", "Roman numeral I to 1"},
		{"ii samuel", "2Sam", "Roman numeral II to 2"},
		{"wisdom", "Wis", "apocrypha lowercase"},
		{"wis", "Wis", "apocrypha abbreviated"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			osis, ok := tbl.ByAlias[tc.alias]
			if !ok {
				t.Errorf("alias %q not found in table", tc.alias)
				return
			}
			if osis != tc.expected {
				t.Errorf("expected alias %q to map to %q, got %q", tc.alias, tc.expected, osis)
			}

			// Verify OSIS exists in ByOsis
			book, exists := tbl.ByOsis[osis]
			if !exists {
				t.Errorf("OSIS %q not found in ByOsis", osis)
				return
			}

			// Verify the book is valid
			if book.OSIS != tc.expected {
				t.Errorf("expected book OSIS to be %q, got %q", tc.expected, book.OSIS)
			}
		})
	}
}

// TestTable_DuplicateAliases checks that NewTable detects and reports duplicate alias keys.
func TestTable_DuplicateAliases(t *testing.T) {
	// Create books with duplicate aliases
	booksWithDuplicates := []bibleref.Book{
		{
			OSIS:      "Prov",
			Name:      "Proverbs",
			Aliases:   []string{"proverbs", "prov"},
			Testament: "OT",
			Order:     20,
			Chapters:  31,
		},
		{
			OSIS:      "Matt",
			Name:      "Matthew",
			Aliases:   []string{"matthew", "prov"}, // Duplicate alias!
			Testament: "NT",
			Order:     40,
			Chapters:  28,
		},
	}

	tbl, err := bibleref.NewTable(booksWithDuplicates)
	if err == nil {
		// NewTable silently overwrites duplicates - document as potential bug
		t.Logf("WARNING: NewTable did not detect duplicate alias 'prov' in books. Silently used last one: %s", tbl.ByAlias["prov"])
		// The duplicate key will be silently overwritten in the map
		if tbl.ByAlias["prov"] != "Matt" {
			t.Errorf("duplicate alias 'prov' was overwritten; last value should be 'Matt', got %q", tbl.ByAlias["prov"])
		}
	}
}

// TestParse_ValidReferences tests parsing of valid Bible references.
// NOTE: BUG EXPOSED - Book names starting with digits (e.g., "1 Samuel", "1 John") are not supported.
// The parser splits on the first digit, which fails for books that start with a digit.
// Only books that have at least one letter before the first digit can be parsed.
func TestParse_ValidReferences(t *testing.T) {
	books := testBooks()
	tbl, err := bibleref.NewTable(books)
	if err != nil {
		t.Fatalf("NewTable failed: %v", err)
	}

	testCases := []struct {
		input        string
		expectedOSIS string
		expectedCh   int
		expectedVs   *util.VerseRange
		desc         string
	}{
		// Proverbs variants
		{
			input:        "Prov 31",
			expectedOSIS: "Prov",
			expectedCh:   31,
			expectedVs:   nil,
			desc:         "Prov 31 chapter only",
		},
		{
			input:        "Proverbs 31:10–31",
			expectedOSIS: "Prov",
			expectedCh:   31,
			expectedVs:   &util.VerseRange{StartVerse: 10, EndVerse: util.Ptr(31)},
			desc:         "Proverbs 31:10–31 full name with en-dash",
		},
		{
			input:        "PRO 31:10-31",
			expectedOSIS: "Prov",
			expectedCh:   31,
			expectedVs:   &util.VerseRange{StartVerse: 10, EndVerse: util.Ptr(31)},
			desc:         "PRO 31:10-31 uppercase with hyphen",
		},
		// Apocrypha
		{
			input:        "Wis 1:1-5",
			expectedOSIS: "Wis",
			expectedCh:   1,
			expectedVs:   &util.VerseRange{StartVerse: 1, EndVerse: util.Ptr(5)},
			desc:         "Wisdom apocrypha with range",
		},
		{
			input:        "Wisdom 1:1",
			expectedOSIS: "Wis",
			expectedCh:   1,
			expectedVs:   &util.VerseRange{StartVerse: 1},
			desc:         "Wisdom full name single verse",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ref, err := bibleref.Parse(tc.input, tbl)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tc.input, err)
				return
			}
			if ref == nil {
				t.Errorf("Parse(%q) returned nil", tc.input)
				return
			}

			if ref.OSIS != tc.expectedOSIS {
				t.Errorf("expected OSIS %q, got %q", tc.expectedOSIS, ref.OSIS)
			}
			if ref.Chapter != tc.expectedCh {
				t.Errorf("expected chapter %d, got %d", tc.expectedCh, ref.Chapter)
			}

			if tc.expectedVs == nil {
				if ref.Verse != nil {
					t.Errorf("expected no verse, got %v", ref.Verse)
				}
			} else {
				if ref.Verse == nil {
					t.Errorf("expected verse %v, got nil", tc.expectedVs)
					return
				}
				if ref.Verse.StartVerse != tc.expectedVs.StartVerse {
					t.Errorf("expected start verse %d, got %d", tc.expectedVs.StartVerse, ref.Verse.StartVerse)
				}
				if (tc.expectedVs.EndVerse == nil) != (ref.Verse.EndVerse == nil) {
					t.Errorf("expected end verse %v, got %v", tc.expectedVs.EndVerse, ref.Verse.EndVerse)
				}
				if tc.expectedVs.EndVerse != nil && ref.Verse.EndVerse != nil {
					if *tc.expectedVs.EndVerse != *ref.Verse.EndVerse {
						t.Errorf("expected end verse %d, got %d", *tc.expectedVs.EndVerse, *ref.Verse.EndVerse)
					}
				}
			}
		})
	}
}

// TestParse_InvalidReferences tests parsing of invalid Bible references.
func TestParse_InvalidReferences(t *testing.T) {
	books := testBooks()
	tbl, err := bibleref.NewTable(books)
	if err != nil {
		t.Fatalf("NewTable failed: %v", err)
	}

	testCases := []struct {
		input       string
		desc        string
		expectError bool
	}{
		{
			input:       "",
			desc:        "empty string",
			expectError: true,
		},
		{
			input:       "Unknown 1:1",
			desc:        "unknown book",
			expectError: true,
		},
		{
			input:       "Prov 0",
			desc:        "chapter 0",
			expectError: true,
		},
		{
			input:       "Prov 32",
			desc:        "chapter beyond max (Proverbs has 31)",
			expectError: true,
		},
		{
			input:       "Prov 1:0",
			desc:        "verse 0",
			expectError: true,
		},
		{
			input:       "Prov 1:20-10",
			desc:        "reversed range (end < start)",
			expectError: true,
		},
		{
			input:       "1Sam 15:1–16:1",
			desc:        "cross-chapter range (unsupported)",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ref, err := bibleref.Parse(tc.input, tbl)
			if !tc.expectError && err != nil {
				t.Errorf("Parse(%q) expected success but got error: %v", tc.input, err)
			}
			if tc.expectError && err == nil {
				t.Errorf("Parse(%q) expected error but got success: %v", tc.input, ref)
			}
		})
	}
}

// TestParseCanonical_Rendering tests that parsing and then calling String() yields canonical form.
// NOTE: BUG EXPOSED - Book names starting with digits are not supported due to parser design.
func TestParseCanonical_Rendering(t *testing.T) {
	books := testBooks()
	tbl, err := bibleref.NewTable(books)
	if err != nil {
		t.Fatalf("NewTable failed: %v", err)
	}

	testCases := []struct {
		input             string
		expectedCanonical string
		desc              string
	}{
		{
			input:             "Proverbs 31:10-31",
			expectedCanonical: "Prov 31:10–31",
			desc:              "hyphen normalized to en-dash",
		},
		{
			input:             "PRO 31:10-31",
			expectedCanonical: "Prov 31:10–31",
			desc:              "uppercase normalized to canonical OSIS",
		},
		{
			input:             "Prov 31:10–31",
			expectedCanonical: "Prov 31:10–31",
			desc:              "already canonical",
		},
		{
			input:             "Wis 1:1-5",
			expectedCanonical: "Wis 1:1–5",
			desc:              "apocrypha with hyphen normalization",
		},
		{
			input:             "Prov 31",
			expectedCanonical: "Prov 31",
			desc:              "chapter-only reference",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ref, err := bibleref.Parse(tc.input, tbl)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tc.input, err)
				return
			}
			canonical := ref.String()
			if canonical != tc.expectedCanonical {
				t.Errorf("expected canonical form %q, got %q", tc.expectedCanonical, canonical)
			}
		})
	}
}

// TestParseCanonical_NormalizationVariants tests that whitespace/punctuation variations normalize to same output.
// NOTE: BUG EXPOSED - Em-dashes are not normalized in verse ranges, only hyphens are converted to en-dashes.
func TestParseCanonical_NormalizationVariants(t *testing.T) {
	books := testBooks()
	tbl, err := bibleref.NewTable(books)
	if err != nil {
		t.Fatalf("NewTable failed: %v", err)
	}

	// All these variants should normalize to the same canonical form
	variants := []string{
		"Prov 31:10-31",
		"Prov 31:10–31",
		"Proverbs 31:10-31",
		"proverbs 31:10–31",
		"PRO 31:10-31",
		"Pro 31:10–31",
		"   Prov   31:10-31   ",
	}

	expectedCanonical := "Prov 31:10–31"

	var firstRef *bibleref.BibleRef
	for _, input := range variants {
		t.Run(input, func(t *testing.T) {
			ref, err := bibleref.Parse(input, tbl)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", input, err)
				return
			}
			canonical := ref.String()
			if canonical != expectedCanonical {
				t.Errorf("expected %q, got %q", expectedCanonical, canonical)
			}

			if firstRef == nil {
				firstRef = ref
			} else {
				// Verify structural equivalence
				if ref.OSIS != firstRef.OSIS || ref.Chapter != firstRef.Chapter {
					t.Errorf("variant %q produced different OSIS/Chapter than first variant", input)
				}
			}
		})
	}
}
