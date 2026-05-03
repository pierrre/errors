// Package errslog provides utilities to manage errors with the slog package.
package errslog

import (
	"log/slog"
)

func WrapAttr(err error, attr slog.Attr) error {
	if err == nil {
		return nil
	}
	return &attrError{
		error: err,
		attr:  attr,
	}
}

type attrError struct {
	error
	attr slog.Attr
}

func (e *attrError) Unwrap() error {
	return e.error
}

func (e *attrError) SlogAttr() slog.Attr {
	return e.attr
}
