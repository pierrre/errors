package errverbose_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func ExampleWrite() {
	err := errbase.New("error")
	buf := new(strings.Builder)
	Write(buf, err)
	s := buf.String()
	fmt.Println(s)
	// Output: error
}

func ExampleString() {
	err := errbase.New("error")
	s := String(err)
	fmt.Println(s)
	// Output: error
}

func ExampleFormatter() {
	err := errbase.New("error")
	f := Formatter(err)
	fmt.Println(f)
	// Output: error
}

type testVerbose struct {
	error
}

func (v *testVerbose) ErrorVerbose() string {
	return "verbose"
}

func (v *testVerbose) Unwrap() error {
	return v.error
}

func TestWrite(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	buf := new(strings.Builder)
	Write(buf, err)
	s := buf.String()
	assert.Equal(t, s, "error\nverbose\n")
}

func TestString(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	s := String(err)
	assert.Equal(t, s, "error\nverbose\n")
}

func TestFormatter(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	f := Formatter(err)
	buf := new(strings.Builder)
	_, _ = fmt.Fprintf(buf, "%v", f)
	s := buf.String()
	assert.Equal(t, s, "error\nverbose\n")
}

func TestNil(t *testing.T) {
	s := String(nil)
	assert.Equal(t, s, "<nil>\n")
}

func TestJoin(t *testing.T) {
	err := &testVerbose{
		error: errors.Join(
			&testVerbose{
				error: errors.Join(
					&testVerbose{
						error: errbase.New("error a"),
					},
					&testVerbose{
						error: errbase.New("error b"),
					},
				),
			},
			&testVerbose{
				error: errors.Join(
					&testVerbose{
						error: errbase.New("error c"),
					},
					&testVerbose{
						error: errbase.New("error d"),
					},
				),
			},
		),
	}
	s := String(err)
	expected := `error a
error b
error c
error d
verbose

Sub error 0: error a
error b
verbose

Sub error 0.0: error a
verbose

Sub error 0.1: error b
verbose

Sub error 1: error c
error d
verbose

Sub error 1.0: error c
verbose

Sub error 1.1: error d
verbose
`
	assert.Equal(t, s, expected)
}

func TestWriteAllocs(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	assert.AllocsPerRun(t, 100, func() {
		Write(io.Discard, err)
	}, 0)
}

func BenchmarkWrite(b *testing.B) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	for i := 0; i < b.N; i++ {
		Write(io.Discard, err)
	}
}
