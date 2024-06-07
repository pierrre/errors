package erriter_test

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/errors/errmsg"
)

func newtestError() error {
	err := errbase.New("error")
	err = errors.Join(err, err)
	err = errmsg.Wrap(err, "test")
	return err
}

func TestIter(t *testing.T) {
	err := newtestError()
	count := 0
	erriter.Iter(err, func(err error) {
		count++
		assert.Error(t, err)
	})
	assert.Equal(t, count, 5)
}

func TestIterAllocs(t *testing.T) {
	err := newtestError()
	assert.AllocsPerRun(t, 100, func() {
		erriter.Iter(err, func(err error) {})
	}, 0)
}

func BenchmarkIter(b *testing.B) {
	err := newtestError()
	for i := 0; i < b.N; i++ {
		erriter.Iter(err, func(err error) {})
	}
}
