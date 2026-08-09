package erriter_test

import (
	std_errors "errors"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/errors/errmsg"
)

func newTestError() error {
	err := errbase.New("error")
	err = errors.Join(err, err)
	err = errmsg.Wrap(err, "test")
	return err
}

func TestAll(t *testing.T) {
	err := newTestError()
	count := 0
	for err := range erriter.All(err) {
		count++
		assert.Error(t, err)
	}
	assert.Equal(t, count, 5)
}

func TestAllStop(t *testing.T) {
	err := newTestError()
	count := 0
	for range erriter.All(err) {
		count++
		if count == 4 {
			break
		}
	}
	assert.Equal(t, count, 4)
}

func TestAllTraversalOrder(t *testing.T) {
	a := errbase.New("a")
	b := errbase.New("b")
	c := errbase.New("c")
	j1 := std_errors.Join(a, b)
	w1 := errmsg.Wrap(j1, "w1")
	w2 := errmsg.Wrap(c, "w2")
	j0 := std_errors.Join(w1, w2)
	err := errmsg.Wrap(j0, "w0")
	var got []error
	for e := range erriter.All(err) {
		got = append(got, e)
	}
	assert.SliceEqual(t, got, []error{err, j0, w1, j1, a, b, w2, c})
}

func TestAllAllocs(t *testing.T) {
	err := newTestError()
	assert.AllocsPerRun(t, 100, func() {
		for range erriter.All(err) {
		}
	}, 0)
}

func BenchmarkAll(b *testing.B) {
	err := newTestError()
	for b.Loop() {
		for range erriter.All(err) {
		}
	}
}
