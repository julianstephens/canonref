package bibleref

import "fmt"

type ErrKind int

const (
	KindParse ErrKind = iota
	KindUnknownBook
	KindInvalidBook
	KindInvalidChapter
	KindInvalidVerse
	KindUnsupportedFormat
)

var (
	ErrBibleRefParseFailed      = fmt.Errorf("parse failed")
	ErrBibleRefValidationFailed = fmt.Errorf("validation failed")
	ErrInvalidOSISCode          = fmt.Errorf("invalid OSIS code")
	ErrInvalidBook              = fmt.Errorf("invalid book")
	ErrInvalidChapter           = fmt.Errorf("invalid chapter")
	ErrInvalidVerse             = fmt.Errorf("invalid verse")
	ErrUnsupportedFormat        = fmt.Errorf("unsupported format")
)

type BibleRefError struct {
	Kind    ErrKind
	Message *string
	Err     error
	Cause   error
}

func (e *BibleRefError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("Bible reference error: %s, err: %v (cause: %v)", *e.Message, e.Err, e.Cause)
	}
	return fmt.Sprintf("Bible reference error: err: %v (cause: %v)", e.Err, e.Cause)
}

func (e *BibleRefError) Unwrap() error {
	return e.Err
}
