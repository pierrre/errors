package errappend_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errmsg"
)

var testSink any

func Example() {
	err := errbase.New("error")
	b := Append(nil, err)
	fmt.Println(string(b))
	// Output: error
}

type testAppendError struct {
	msg string
}

func (e *testAppendError) Error() string {
	return "not used"
}

func (e *testAppendError) ErrorAppend(b []byte) []byte {
	return append(b, e.msg...)
}

func TestAppend(t *testing.T) {
	err := errbase.New("error")
	b := Append(nil, err)
	assert.Equal(t, string(b), "error")
}

func TestAppendPrefix(t *testing.T) {
	err := errbase.New("error")
	b := Append([]byte("prefix: "), err)
	assert.Equal(t, string(b), "prefix: error")
}

func TestAppendInterface(t *testing.T) {
	err := &testAppendError{
		msg: "error",
	}
	b := Append(nil, err)
	assert.Equal(t, string(b), "error")
}

func TestAppendInterfaceWrapped(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	b := Append(nil, err)
	assert.Equal(t, string(b), "msg: error")
}

func TestAppendNil(t *testing.T) {
	b := []byte("test")
	res := Append(b, nil)
	assert.Equal(t, string(res), string(b))
}

func TestString(t *testing.T) {
	err := &testAppendError{
		msg: "error",
	}
	assert.Equal(t, String(err), "error")
}

func TestStringWrapped(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	i, _ := assert.ErrorAsType[Interface](t, err)
	assert.Equal(t, String(i), "msg: error")
}

func TestStringNil(t *testing.T) {
	assert.Equal(t, String(nil), "")
}

func TestAppendAllocs(t *testing.T) {
	err := errbase.New("error")
	var res []byte
	assert.AllocsPerRun(t, 100, func() {
		res = Append(nil, err)
	}, 1)
	testSink = res
}

func TestStringAllocs(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	i, _ := assert.ErrorAsType[Interface](t, err)
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = String(i)
	}, 1)
	testSink = res
}

func BenchmarkAppend(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Append(nil, err)
	}
}

func BenchmarkString(b *testing.B) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	i, _ := assert.ErrorAsType[Interface](b, err)
	for b.Loop() {
		_ = String(i)
	}
}
