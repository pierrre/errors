// Package errbase provides a simple error type.
package errbase

type base struct {
	msg string
}

// New creates a new error with the given message.
//
// It can be used as a "sentinel" error.
func New(msg string) error {
	return &base{
		msg: msg,
	}
}

func (err *base) Error() string {
	return err.msg
}
