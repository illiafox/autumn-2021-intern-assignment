package public

import "errors"

type InternalError struct {
	err error
}

func (i InternalError) Error() string {
	return i.err.Error()
}

func NewInternal(err error) error {
	return InternalError{err}
}

// ErrInternal is used to verify NewInternal error
//
// public.AsInternal(err)
var errCheck = &InternalError{}

func AsInternal(err error) bool {
	return errors.As(err, errCheck)
}
