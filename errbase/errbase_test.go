package errbase_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func Example() {
	err := New("error")
	s := err.Error()
	fmt.Println(s)
	// Output: error
}

func Test(t *testing.T) {
	err := Newf("error %d", 1)
	assert.ErrorEqual(t, err, "error 1")
}

func TestNewAllocs(t *testing.T) {
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = New("error")
	}, 1)
	runtime.KeepAlive(res)
}

func TestNewfAllocs(t *testing.T) {
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Newf("error %d", 1)
	}, 2)
	runtime.KeepAlive(res)
}

func TestErrorAllocs(t *testing.T) {
	err := New("error")
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 0)
	runtime.KeepAlive(res)
}

func BenchmarkNew(b *testing.B) {
	var res error
	for i := 0; i < b.N; i++ {
		res = New("error")
	}
	runtime.KeepAlive(res)
}

func BenchmarkNewf(b *testing.B) {
	var res error
	for i := 0; i < b.N; i++ {
		res = Newf("error %d", 1)
	}
	runtime.KeepAlive(res)
}

func BenchmarkError(b *testing.B) {
	err := New("error")
	var res string
	for i := 0; i < b.N; i++ {
		res = err.Error()
	}
	runtime.KeepAlive(res)
}
