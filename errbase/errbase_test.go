package errbase_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors/errbase"
)

var testSink any

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
	testSink = res
}

func TestNewfAllocs(t *testing.T) {
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Newf("error %d", 1)
	}, 2)
	testSink = res
}

func TestErrorAllocs(t *testing.T) {
	err := New("error")
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 0)
	testSink = res
}

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_ = New("error")
	}
}

func BenchmarkNewf(b *testing.B) {
	for b.Loop() {
		_ = Newf("error %d", 1)
	}
}

func BenchmarkError(b *testing.B) {
	err := New("error")
	for b.Loop() {
		_ = err.Error()
	}
}
