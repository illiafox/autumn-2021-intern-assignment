package public

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
// errors.As(err, public.ErrInternal)
var ErrInternal = &InternalError{}
