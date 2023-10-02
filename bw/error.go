package bw

type BWError struct {
	err error
	msg string
}

func NewBWError(err error, message string) *BWError {
	return &BWError{
		err: err,
		msg: message,
	}
}

func (e *BWError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func (e *BWError) Unwrap() error {
	return e.err
}
