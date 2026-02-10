package rbref

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/julianstephens/go-utils/helpers"

	"github.com/julianstephens/canonref/util"
)

type RbSectionKind uint8

const (
	RbPrologue RbSectionKind = iota
	RbChapter
)

type RbRef struct {
	Kind       RbSectionKind    `json:"kind"`
	ChapterNum *int             `json:"chapter_num,omitempty"`
	Verse      *util.VerseRange `json:"verse,omitempty"`
}

func NewRbRef(rbStr string) (*RbRef, error) {
	ref, err := parseRbRef(rbStr)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// validate checks that the RbRef has valid values based on its kind (prologue or chapter).
// For prologue references, chapter number must be nil and verse range must be present and valid.
// For chapter references, chapter number must be present and between 1 and 73, and verse range (if present) must be valid.
func (r *RbRef) validate() error {
	validateVerse := func(v *util.VerseRange) error {
		if v.StartVerse < 1 {
			return &RbRefError{
				Err:     ErrRbRefValidationFailed,
				Message: util.Ptr("RB reference start verse must be positive and at least 1"),
			}
		}
		if v.EndVerse != nil {
			if *v.EndVerse < v.StartVerse {
				return &RbRefError{
					Err:     ErrRbRefValidationFailed,
					Message: util.Ptr("RB reference end verse cannot be less than start verse"),
				}
			}
		}
		return nil
	}

	switch r.Kind {
	case RbPrologue:
		if r.ChapterNum != nil {
			return &RbRefError{
				Err:     ErrRbRefValidationFailed,
				Message: util.Ptr("prologue RB reference cannot have chapter number"),
			}
		}
		if r.Verse == nil {
			return &RbRefError{
				Err:     ErrRbRefValidationFailed,
				Message: util.Ptr("prologue RB reference must have verse range"),
			}
		}
		if err := validateVerse(r.Verse); err != nil {
			return err
		}
	case RbChapter:
		if r.ChapterNum == nil {
			return &RbRefError{
				Err:     ErrRbRefValidationFailed,
				Message: util.Ptr("chapter RB reference must have chapter number"),
			}
		}
		if *r.ChapterNum < 1 || *r.ChapterNum > 73 {
			return &RbRefError{
				Err:     ErrRbRefValidationFailed,
				Message: util.Ptr("chapter RB reference chapter number must be positive and between 1 and 73"),
			}
		}
		if r.Verse != nil {
			if err := validateVerse(r.Verse); err != nil {
				return err
			}
		}
	default:
		return &RbRefError{
			Err:     ErrRbRefValidationFailed,
			Message: util.Ptr("invalid RB reference kind"),
		}
	}
	return nil
}

func (r *RbRef) String() string {
	var versePart string
	if r.Verse != nil {
		versePart = fmt.Sprintf("%d", r.Verse.StartVerse)
		if r.Verse.EndVerse != nil {
			versePart += fmt.Sprintf("%s%d", util.EnDash, *r.Verse.EndVerse)
		}
	}

	switch r.Kind {
	case RbPrologue:
		return fmt.Sprintf("RB Prol. %s", versePart)
	case RbChapter:
		if r.ChapterNum != nil {
			return fmt.Sprintf("RB %d%s", *r.ChapterNum, util.If(versePart != "", fmt.Sprintf(".%s", versePart), ""))
		}
		return "RB (invalid chapter reference)"
	default:
		return "RB (invalid reference)"
	}
}

func parseRbRef(rbStr string) (*RbRef, error) {
	var ref *RbRef

	parts := strings.Split(rbStr, " ")
	if len(parts) < 2 {
		return nil, &RbRefError{
			Err: ErrRbRefParseFailed,
			Message: util.Ptr(fmt.Sprintf(
				"invalid RB reference format: %s",
				rbStr,
			)),
		}
	}

	work := parts[0]
	if work != "RB" {
		return nil, &RbRefError{
			Err: ErrRbRefParseFailed,
			Message: util.Ptr(fmt.Sprintf(
				"invalid RB reference format: %s",
				rbStr,
			)),
		}
	}

	// chapter reference format: "RB 2.1-5" or "RB 2.1"
	if len(parts) == 2 {
		if strings.HasPrefix(parts[1], "Prol.") || strings.HasPrefix(parts[1], "Prol") {
			return nil, &RbRefError{
				Err: ErrRbRefParseFailed,
				Message: util.Ptr(fmt.Sprintf(
					"prologue RB references must include chapter/verse numbers: %s",
					rbStr,
				)),
			}
		}
		chapterRefParts := strings.Split(parts[1], ".")
		if len(chapterRefParts) > 2 {
			return nil, &RbRefError{
				Err: ErrRbRefParseFailed,
				Message: util.Ptr(fmt.Sprintf(
					"invalid RB chapter reference format: %s",
					rbStr,
				)),
			}
		}

		chapterNum, err := parsePositiveInt(chapterRefParts[0])
		if err != nil {
			return nil, err
		}

		if len(chapterRefParts) == 1 {
			ref = &RbRef{
				Kind:       RbChapter,
				ChapterNum: &chapterNum,
				Verse:      nil,
			}
		} else {
			verseRange, err := parseVerseRange(helpers.If(len(chapterRefParts) == 2, chapterRefParts[1], ""))
			if err != nil {
				return nil, err
			}

			ref = &RbRef{
				Kind:       RbChapter,
				ChapterNum: &chapterNum,
				Verse:      verseRange,
			}
		}
	}

	// prologue reference format: "RB Prol. 1-5" or "RB Prol. 1" or "RB Prol 1-5" or "RB Prol 1"
	if len(parts) == 3 {
		if parts[1] != "Prol." && parts[1] != "Prol" {
			return nil, &RbRefError{
				Err: ErrRbRefParseFailed,
				Message: util.Ptr(fmt.Sprintf(
					"invalid RB reference format: %s",
					rbStr,
				)),
			}
		}

		verseRange, err := parseVerseRange(parts[2])
		if err != nil {
			return nil, err
		}
		ref = &RbRef{
			Kind:  RbPrologue,
			Verse: verseRange,
		}
	}

	if ref == nil {
		return nil, &RbRefError{
			Err: ErrRbRefParseFailed,
			Message: util.Ptr(fmt.Sprintf(
				"invalid RB reference format: %s",
				rbStr,
			)),
		}
	}

	if err := ref.validate(); err != nil {
		return nil, err
	}

	return ref, nil
}

func parseVerseRange(verseRangeStr string) (*util.VerseRange, error) {
	if verseRangeStr == "" {
		return nil, nil
	}

	normalized := strings.ReplaceAll(verseRangeStr, "-", "–")
	normalized = strings.ReplaceAll(normalized, "—", "–")
	parts := strings.Split(normalized, "–")

	if len(parts) == 0 || len(parts) > 2 {
		return nil, &RbRefError{
			Err: ErrRbRefParseFailed,
			Message: util.Ptr(fmt.Sprintf(
				"invalid RB verse range format: %s",
				verseRangeStr,
			)),
		}
	}

	startVerse, err := parsePositiveInt(parts[0])
	if err != nil {
		return nil, err
	}

	var endVerse *int
	if len(parts) == 2 {
		ev, err := parsePositiveInt(parts[1])
		if err != nil {
			return nil, err
		}
		endVerse = &ev
	}

	if endVerse != nil {
		if *endVerse < startVerse {
			return nil, &RbRefError{
				Err: ErrRbRefValidationFailed,
				Message: util.Ptr(fmt.Sprintf(
					"end verse cannot be less than start verse: %s",
					verseRangeStr,
				)),
			}
		}

		if *endVerse == startVerse {
			endVerse = nil
		}
	}

	return &util.VerseRange{
		StartVerse: startVerse,
		EndVerse:   endVerse,
	}, nil
}

func parsePositiveInt(s string) (int, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, &RbRefError{
			Err: ErrRbRefParseFailed,
			Message: util.Ptr(fmt.Sprintf(
				"invalid positive integer: %s",
				s,
			)),
			Cause: err,
		}
	}
	if n <= 0 {
		return 0, &RbRefError{
			Err: ErrRbRefParseFailed,
			Message: util.Ptr(fmt.Sprintf(
				"integer must be positive: %s",
				s,
			)),
		}
	}
	return int(n), nil
}
