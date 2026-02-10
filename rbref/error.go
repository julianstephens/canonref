package rbref

import "fmt"

var (
	ErrRbRefValidationFailed = fmt.Errorf("validation failed")
	ErrRbRefParseFailed      = fmt.Errorf("parse failed")
)

type RbRefError struct {
	Message *string
	Err     error
	Cause   error
}

func (e *RbRefError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("rb reference error: %s, err: %v (cause: %v)", *e.Message, e.Err, e.Cause)
	}
	return fmt.Sprintf("rb reference error: err: %v (cause: %v)", e.Err, e.Cause)
}

func (e *RbRefError) Unwrap() error {
	return e.Err
}
