package erriter_test

import (
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

func TestAllAllocs(t *testing.T) {
	err := newTestError()
	assert.AllocsPerRun(t, 100, func() {
		for range erriter.All(err) {
		}
	}, 0)
}

func TestFirstKeys(t *testing.T) {
	seq := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
		yield("a", 3)
		yield("c", 4)
	}
	m := erriter.FirstKeys(seq)
	assert.MapEqual(t, m, map[string]int{
		"a": 1,
		"b": 2,
		"c": 4,
	})
}

func TestFirstKeysEmpty(t *testing.T) {
	seq := func(yield func(string, int) bool) {}
	m := erriter.FirstKeys(seq)
	assert.MapEmpty(t, m)
}

func BenchmarkAll(b *testing.B) {
	err := newTestError()
	for b.Loop() {
		for range erriter.All(err) {
		}
	}
}
