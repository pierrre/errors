//go:build go1.26

package errors

import (
	std_errors "errors"
)

// AsType is an alias for [std_errors.AsType].
func AsType[E error](err error) (E, bool) {
	return std_errors.AsType[E](err)
}
