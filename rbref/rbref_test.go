package rbref_test

import (
	"testing"

	"github.com/julianstephens/canonref/rbref"
	"github.com/julianstephens/canonref/util"
)

func TestNewRbRef_ValidCases(t *testing.T) {
	ok := []string{
		"RB Prol. 1",
		"RB Prol. 1–7",
		"RB Prol 1", // missing dot accepted and normalized
		"RB 48.1",
		"RB 48.1–9",
		"RB 4",
		"RB 4.72-74", // hyphen accepted and normalized to en-dash
		"RB 48.1—9",  // em-dash accepted and normalized to en-dash
		"RB 4.72–74", // en-dash accepted
	}

	for _, tc := range ok {
		t.Run(tc, func(t *testing.T) {
			ref, err := rbref.NewRbRef(tc)
			if err != nil {
				t.Errorf("Expected valid RB reference for %q, got error: %v", tc, err)
			}
			if ref == nil {
				t.Errorf("Expected non-nil RbRef for %q", tc)
			}
		})
	}
}

func TestNewRbRef_InvalidCases(t *testing.T) {
	bad := []string{
		"",
		"RB 0.1",       // chapter 0 invalid
		"RB 74.1",      // chapter too high
		"RB 48.0",      // verse 0 invalid
		"RB 48.9–1",    // reverse range
		"RB Prol. 3–2", // reverse range
		"RB foo",
	}

	for _, tc := range bad {
		t.Run(tc, func(t *testing.T) {
			ref, err := rbref.NewRbRef(tc)
			if err == nil {
				t.Errorf("Expected error for invalid RB reference %q, got nil. Result: %v", tc, ref)
			}
			if ref != nil {
				t.Errorf("Expected nil RbRef for invalid reference %q, got: %v", tc, ref)
			}
		})
	}
}

func TestRbRef_PrologueReferences(t *testing.T) {
	testCases := []struct {
		input       string
		startVerse  int
		endVerse    *int
		description string
	}{
		{
			input:       "RB Prol. 1",
			startVerse:  1,
			endVerse:    nil,
			description: "single verse prologue",
		},
		{
			input:       "RB Prol. 1–7",
			startVerse:  1,
			endVerse:    util.Ptr(7),
			description: "verse range prologue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ref, err := rbref.NewRbRef(tc.input)
			if err != nil {
				t.Fatalf("NewRbRef failed: %v", err)
			}

			if ref.Kind != rbref.RbPrologue {
				t.Errorf("Expected kind RbPrologue, got %v", ref.Kind)
			}

			if ref.ChapterNum != nil {
				t.Errorf("Expected nil chapter number for prologue, got %v", *ref.ChapterNum)
			}

			if ref.Verse == nil {
				t.Fatal("Expected non-nil verse range")
			}

			if ref.Verse.StartVerse != tc.startVerse {
				t.Errorf("Expected start verse %d, got %d", tc.startVerse, ref.Verse.StartVerse)
			}

			if tc.endVerse == nil && ref.Verse.EndVerse != nil {
				t.Errorf("Expected nil end verse, got %d", *ref.Verse.EndVerse)
			}

			if tc.endVerse != nil {
				if ref.Verse.EndVerse == nil {
					t.Errorf("Expected end verse %d, got nil", *tc.endVerse)
				} else if *ref.Verse.EndVerse != *tc.endVerse {
					t.Errorf("Expected end verse %d, got %d", *tc.endVerse, *ref.Verse.EndVerse)
				}
			}
		})
	}
}

func TestRbRef_ChapterReferences(t *testing.T) {
	testCases := []struct {
		input       string
		chapter     int
		startVerse  *int
		endVerse    *int
		description string
	}{
		{
			input:       "RB 4",
			chapter:     4,
			startVerse:  nil,
			endVerse:    nil,
			description: "chapter only",
		},
		{
			input:       "RB 48.1",
			chapter:     48,
			startVerse:  util.Ptr(1),
			endVerse:    nil,
			description: "chapter with single verse",
		},
		{
			input:       "RB 48.1–9",
			chapter:     48,
			startVerse:  util.Ptr(1),
			endVerse:    util.Ptr(9),
			description: "chapter with verse range",
		},
		{
			input:       "RB 4.72-74",
			chapter:     4,
			startVerse:  util.Ptr(72),
			endVerse:    util.Ptr(74),
			description: "hyphen normalized to en dash",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ref, err := rbref.NewRbRef(tc.input)
			if err != nil {
				t.Fatalf("NewRbRef failed: %v", err)
			}

			if ref.Kind != rbref.RbChapter {
				t.Errorf("Expected kind RbChapter, got %v", ref.Kind)
			}

			if ref.ChapterNum == nil {
				t.Fatal("Expected non-nil chapter number")
			}

			if *ref.ChapterNum != tc.chapter {
				t.Errorf("Expected chapter %d, got %d", tc.chapter, *ref.ChapterNum)
			}

			if tc.startVerse == nil {
				if ref.Verse != nil {
					t.Errorf("Expected nil verse, got %+v", ref.Verse)
				}
			} else {
				if ref.Verse == nil {
					t.Fatal("Expected non-nil verse")
				}

				if ref.Verse.StartVerse != *tc.startVerse {
					t.Errorf("Expected start verse %d, got %d", *tc.startVerse, ref.Verse.StartVerse)
				}

				if tc.endVerse == nil {
					if ref.Verse.EndVerse != nil {
						t.Errorf("Expected nil end verse, got %d", *ref.Verse.EndVerse)
					}
				} else {
					if ref.Verse.EndVerse == nil {
						t.Errorf("Expected end verse %d, got nil", *tc.endVerse)
					} else if *ref.Verse.EndVerse != *tc.endVerse {
						t.Errorf("Expected end verse %d, got %d", *tc.endVerse, *ref.Verse.EndVerse)
					}
				}
			}
		})
	}
}
