package errslog

import (
	"fmt"
	"log/slog"

	"github.com/pierrre/go-libs/bufpool"
)

func WrapAttrs(err error, attrs ...slog.Attr) error {
	if err == nil {
		return nil
	}
	if len(attrs) == 0 {
		return err
	}
	buf := bufferPool.Get()
	defer bufferPool.Put(buf)
	for i, attr := range attrs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(attr.Key)
		buf.WriteString("=")
		buf.WriteString(attr.Value.String())
	}
	buf.WriteString(": ")
	buf.WriteString(err.Error())
	return &attrsError{
		error: err,
		msg:   buf.String(),
		attrs: attrs,
	}
}

type attrsError struct {
	error
	msg   string
	attrs []slog.Attr
}

func (err *attrsError) Unwrap() error {
	return err.error
}

func (err *attrsError) Error() string {
	return err.msg
}

func WrapLevel(err error, level slog.Level) error {
	if err == nil {
		return nil
	}
	return &levelError{
		error: err,
		msg:   fmt.Sprintf("level=%v: %v", level, err),
		level: level,
	}
}

type levelError struct {
	error
	msg   string
	level slog.Level
}

func (err *levelError) Unwrap() error {
	return err.error
}

func (err *levelError) Error() string {
	return err.msg
}

var bufferPool = bufpool.Pool{}
