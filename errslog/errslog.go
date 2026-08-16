// Package errslog provides utilities to use errors with the slog package:
// wrap an error with attributes ([WrapAttrs]) and a level ([WrapLevel]), and log it with [Log] or [LoggerLog].
package errslog

import (
	"context"
	"errors"
	"iter"
	"log/slog"
	"slices"
	"strconv"

	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/syncutil"
)

// WrapAttrs wraps an error with attributes ([slog.Attr]).
// It returns nil if err is nil.
// If attrs is empty, it returns err.
// The error keeps a reference to the attrs slice and does not copy it; the caller must not reuse or modify it after the call.
// Use [AllAttrs] or [GetAttrs] to retrieve the attributes.
func WrapAttrs(err error, attrs ...slog.Attr) error {
	if err == nil {
		return nil
	}
	if len(attrs) == 0 {
		return err
	}
	return &attrError{
		error: err,
		attrs: attrs,
	}
}

type attrError struct {
	error
	attrs []slog.Attr
}

func (e *attrError) Unwrap() error {
	return e.error
}

func (e *attrError) SlogAttrs() []slog.Attr {
	return e.attrs
}

func (e *attrError) Error() string {
	return errappend.String(e)
}

func (e *attrError) ErrorAppend(b []byte) []byte {
	b = appendAttributes(b, e.attrs)
	b = append(b, ": "...)
	b = errappend.Append(b, e.error)
	return b
}

func appendAttributes(b []byte, attrs []slog.Attr) []byte {
	for i, attr := range attrs {
		if i > 0 {
			b = append(b, ", "...)
		}
		b = appendAttribute(b, attr)
	}
	return b
}

func appendAttribute(b []byte, attr slog.Attr) []byte {
	b = append(b, attr.Key...)
	b = append(b, '=')
	b = appendValue(b, attr.Value)
	return b
}

func appendValue(b []byte, v slog.Value) []byte {
	switch v.Kind() { //nolint:exhaustive // Some types can be optimized, the rest can be handled by String().
	case slog.KindString:
		return strconv.AppendQuote(b, v.String())
	case slog.KindBool:
		return append(b, strconv.FormatBool(v.Bool())...)
	case slog.KindFloat64:
		return strconv.AppendFloat(b, v.Float64(), 'g', -1, 64)
	case slog.KindInt64:
		return strconv.AppendInt(b, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(b, v.Uint64(), 10)
	case slog.KindGroup:
		return appendGroup(b, v.Group())
	default:
		return append(b, v.String()...)
	}
}

func appendGroup(b []byte, attrs []slog.Attr) []byte {
	b = append(b, '[')
	b = appendAttributes(b, attrs)
	b = append(b, ']')
	return b
}

// AllAttrs returns an [iter.Seq] that iterates over all the attributes of an error tree recursively.
// The order is from the outermost error to the innermost (see [erriter.All]).
// Duplicates are not removed; use [GetAttrs] to get unique attributes.
func AllAttrs(err error) iter.Seq[slog.Attr] {
	return func(yield func(slog.Attr) bool) {
		for err := range erriter.All(err) {
			erra, ok := err.(interface{ SlogAttrs() []slog.Attr })
			if !ok {
				continue
			}
			for _, attr := range erra.SlogAttrs() {
				if !yield(attr) {
					return
				}
			}
		}
	}
}

// GetAttrs returns all the attributes of an error tree, recursively, without duplicates.
// For each key, only the first attribute is kept (see [AllAttrs]).
func GetAttrs(err error) []slog.Attr {
	attrs := getAttrs(err)
	defer releaseAttrsToPool(attrs)
	var res []slog.Attr
	if len(attrs) > 0 {
		res = make([]slog.Attr, len(attrs))
		copy(res, attrs)
	}
	return res
}

func getAttrs(err error) []slog.Attr {
	attrs := getAttrsFromPool()
	for attr := range AllAttrs(err) {
		if !slices.ContainsFunc(attrs, func(a slog.Attr) bool { // The slices is usually small, so it's OK.
			return a.Key == attr.Key
		}) {
			attrs = append(attrs, attr)
		}
	}
	return attrs
}

var attrsPool = &syncutil.ValuePool[[]slog.Attr]{}

func getAttrsFromPool() []slog.Attr {
	return attrsPool.Get()[:0]
}

func releaseAttrsToPool(attrs []slog.Attr) {
	clear(attrs[:cap(attrs)])
	attrs = attrs[:0]
	attrsPool.Put(attrs)
}

// WrapLevel wraps an error with a [slog.Level] (see [GetLevel]).
// It returns nil if err is nil.
func WrapLevel(err error, l slog.Level) error {
	if err == nil {
		return nil
	}
	return &levelError{
		error: err,
		level: l,
	}
}

type levelError struct {
	error
	level slog.Level
}

func (e *levelError) Unwrap() error {
	return e.error
}

func (e *levelError) ErrorAppend(b []byte) []byte {
	return errappend.Append(b, e.error)
}

func (e *levelError) SlogLevel() slog.Level {
	return e.level
}

func (e *levelError) ErrorVerboseAppend(b []byte) []byte {
	b = append(b, "slog level "...)
	b = append(b, e.level.String()...)
	return b
}

// GetLevel returns the [slog.Level] associated with an error, if it was wrapped with [WrapLevel].
// The ok boolean indicates whether a level is associated with the error.
func GetLevel(err error) (l slog.Level, ok bool) {
	lerr, ok := errors.AsType[interface {
		error
		SlogLevel() slog.Level
	}](err)
	if ok {
		l = lerr.SlogLevel()
	}
	return l, ok
}

// Log calls [LoggerLog] with [slog.Default].
func Log(ctx context.Context, err error, attrs ...slog.Attr) {
	LoggerLog(ctx, nil, err, attrs...)
}

// LoggerLog logs an error with a logger.
// It does nothing if err is nil.
// It uses the [slog.Default] logger if logger is nil.
// It logs at the level [slog.LevelError], with err.Error() as the message, and the error tree attributes ([GetAttrs]) as additional attributes.
func LoggerLog(ctx context.Context, logger *slog.Logger, err error, attrs ...slog.Attr) {
	if err == nil {
		return
	}
	if logger == nil {
		logger = slog.Default()
	}
	level, ok := GetLevel(err)
	if !ok {
		level = slog.LevelError
	}
	if !logger.Enabled(ctx, level) {
		return
	}
	errAttrs := getAttrs(err)
	defer releaseAttrsToPool(errAttrs)
	if len(errAttrs) > 0 {
		if len(attrs) > 0 {
			newAttrs := getAttrsFromPool()
			newAttrs = slices.Grow(newAttrs, len(attrs)+len(errAttrs))
			newAttrs = append(newAttrs, attrs...)
			newAttrs = append(newAttrs, errAttrs...)
			defer releaseAttrsToPool(newAttrs)
			attrs = newAttrs
		} else {
			attrs = errAttrs
		}
	}
	logger.LogAttrs(ctx, level, err.Error(), attrs...)
}
