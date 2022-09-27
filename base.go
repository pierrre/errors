package errors

type base struct {
	s string
}

func newBase(s string) error {
	return &base{
		s: s,
	}
}

func (err *base) Error() string {
	return err.s
}
