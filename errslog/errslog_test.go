package errslog_test

import (
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errslog"
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/go-libs/bytesutil"
)

var testSink any

func Example() {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123))
	attrs := GetAttrs(err)
	fmt.Println(attrs[0])
	// Output: int=123
}

func TestAttrs(t *testing.T) {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123))
	err = WrapAttrs(err, slog.String("string", "test"))
	attrs := GetAttrs(err)
	expected := []slog.Attr{
		slog.String("string", "test"),
		slog.Int("int", 123),
	}
	assert.DeepEqual(t, attrs, expected)
}

func TestWrapAttrsNil(t *testing.T) {
	err := WrapAttrs(nil, slog.Int("int", 123))
	assert.NoError(t, err)
}

func TestWrapAttrsEmpty(t *testing.T) {
	err := errbase.New("error")
	err2 := WrapAttrs(err)
	assert.Equal(t, err, err2)
}

func TestWrapAttrsAlloc(t *testing.T) {
	err := errbase.New("error")
	attrs := []slog.Attr{
		slog.Int("int", 123),
		slog.String("string", "test"),
	}
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = WrapAttrs(err, attrs...)
	}, 1)
	testSink = res
}

func BenchmarkWrapAttrs(b *testing.B) {
	err := errbase.New("error")
	attrs := []slog.Attr{
		slog.Int("int", 123),
		slog.String("string", "test"),
	}
	var e error
	for b.Loop() {
		e = WrapAttrs(err, attrs...)
	}
	testSink = e
}

func TestAttrsError(t *testing.T) {
	err := errbase.New("error")
	err = WrapAttrs(err,
		slog.String("string", "test"),
		slog.Bool("bool", true),
		slog.Float64("float64", 123.456),
		slog.Int64("int64", 123),
		slog.Uint64("uint64", 123),
		slog.Duration("duration", 1*time.Second),
		slog.GroupAttrs("group",
			slog.String("foo", "bar"),
			slog.String("aaa", "zzz"),
		),
	)
	assert.ErrorEqual(t, err, "string=\"test\", bool=true, float64=123.456, int64=123, uint64=123, duration=1s, group=[foo=\"bar\", aaa=\"zzz\"]: error")
}

func TestAttrsErrorAllocs(t *testing.T) {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 1)
	testSink = res
}

func BenchmarkAttrsError(b *testing.B) {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	var s string
	for b.Loop() {
		s = err.Error()
	}
	testSink = s
}

func TestAllAttrsInterrupt(t *testing.T) {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	for range AllAttrs(err) {
		break
	}
}

func TestGetAttrsAllocs(t *testing.T) {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	var res []slog.Attr
	assert.AllocsPerRun(t, 100, func() {
		res = GetAttrs(err)
	}, 1)
	testSink = res
}

func BenchmarkGetAttrs(b *testing.B) {
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	var attrs []slog.Attr
	for b.Loop() {
		attrs = GetAttrs(err)
	}
	testSink = attrs
}

func TestLevel(t *testing.T) {
	err := errbase.New("error")
	err = WrapLevel(err, slog.LevelDebug)
	l, ok := GetLevel(err)
	assert.True(t, ok)
	assert.Equal(t, l, slog.LevelDebug)
	assert.Error(t, errors.Unwrap(err))
}

func TestLevelAppend(t *testing.T) {
	err := errbase.New("error")
	err = WrapLevel(err, slog.LevelDebug)
	_, _ = assert.ErrorAsType[errappend.Interface](t, err)
	b := errappend.Append(nil, err)
	assert.Equal(t, string(b), "error")
}

func TestWrapLevelNil(t *testing.T) {
	err := WrapLevel(nil, slog.LevelDebug)
	assert.NoError(t, err)
}

func TestWrapLevelAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = WrapLevel(err, slog.LevelDebug)
	}, 1)
	testSink = res
}

func BenchmarkWrapLevel(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = WrapLevel(err, slog.LevelDebug)
	}
	testSink = res
}

func TestLevelVerbose(t *testing.T) {
	err := errbase.New("error")
	err = WrapLevel(err, slog.LevelDebug)
	v, _ := assert.ErrorAsType[errverbose.AppendInterface](t, err)
	b := v.ErrorVerboseAppend(nil)
	assert.Equal(t, string(b), "slog level DEBUG")
}

func TestLevelVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = WrapLevel(err, slog.LevelDebug)
	v, _ := assert.ErrorAsType[errverbose.AppendInterface](t, err)
	var b []byte
	assert.AllocsPerRun(t, 100, func() {
		b = v.ErrorVerboseAppend(b)
		b = b[:0]
	}, 0)
	testSink = b
}

func BenchmarkLevelVerbose(b *testing.B) {
	err := errbase.New("error")
	err = WrapLevel(err, slog.LevelDebug)
	v, _ := assert.ErrorAsType[errverbose.AppendInterface](b, err)
	var buf []byte
	for b.Loop() {
		buf = v.ErrorVerboseAppend(buf)
		buf = buf[:0]
	}
	testSink = buf
}

func TestLog(t *testing.T) {
	ctx := t.Context()
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	previousLogger := slog.Default()
	defer slog.SetDefault(previousLogger)
	bw := new(bytesutil.Writer)
	logger := slog.New(slog.NewTextHandler(bw, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(now)
			}
			return a
		},
	}))
	slog.SetDefault(logger)
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	Log(ctx, err)
	Log(ctx, err, slog.String("foo", "bar"))
	expected := "time=2026-01-01T00:00:00.000Z level=ERROR msg=\"int=123, string=\\\"test\\\": error\" int=123 string=test\n" +
		"time=2026-01-01T00:00:00.000Z level=ERROR msg=\"int=123, string=\\\"test\\\": error\" foo=bar int=123 string=test\n"
	assert.Equal(t, bw.String(), expected)
}

func TestLoggerLogErrorNil(t *testing.T) {
	ctx := t.Context()
	bw := new(bytesutil.Writer)
	logger := slog.New(slog.NewTextHandler(bw, nil))
	LoggerLog(ctx, logger, nil)
	assert.Equal(t, bw.String(), "")
}

func TestLoggerLogLevelNotEnabled(t *testing.T) {
	ctx := t.Context()
	bw := new(bytesutil.Writer)
	logger := slog.New(slog.NewTextHandler(bw, nil))
	err := errbase.New("error")
	err = WrapLevel(err, slog.LevelDebug)
	LoggerLog(ctx, logger, err)
	assert.Equal(t, bw.String(), "")
}

func BenchmarkLoggerLog(b *testing.B) {
	ctx := b.Context()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil)) //nolint:sloglint // We can't use DiscardHandler, because we want to write the log.
	err := errbase.New("error")
	err = WrapAttrs(err, slog.Int("int", 123), slog.String("string", "test"))
	attrs := []slog.Attr{slog.String("foo", "bar")}
	for b.Loop() {
		LoggerLog(ctx, logger, err, attrs...)
	}
}
