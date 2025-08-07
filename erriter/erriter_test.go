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

func TestAll(t *testing.T) {
	err := newtestError()
	count := 0
	for err := range erriter.All(err) {
		count++
		assert.Error(t, err)
	}
	assert.Equal(t, count, 5)
}

func TestAllStop(t *testing.T) {
	err := newtestError()
	count := 0
	for range erriter.All(err) {
		count++
		if count == 4 {
			break
		}
	}
	assert.Equal(t, count, 4)
}

func TestAllAllocs(t *testing.T) {
	err := newtestError()
	assert.AllocsPerRun(t, 100, func() {
		for range erriter.All(err) {
		}
	}, 0)
}

func BenchmarkAll(b *testing.B) {
	err := newtestError()
	for b.Loop() {
		for range erriter.All(err) {
		}
	}
}
