package errors

// Wrap adds a message to an error, and optionnally add a stack if it doesn't have one.
//
// See Message() and Stack() for more information.
func Wrap(err error, msg string) error {
	err = Message(err, msg)
	err = ensureStack(err, 2)
	return err
}
