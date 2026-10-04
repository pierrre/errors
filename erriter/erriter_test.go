package erriter_test

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/erriter"
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
	for err := range All(err) {
		count++
		assert.Error(t, err)
	}
	assert.Equal(t, count, 5)
}

func TestAllStop(t *testing.T) {
	err := newTestError()
	count := 0
	for range All(err) {
		count++
		if count == 4 {
			break
		}
	}
	assert.Equal(t, count, 4)
}

type cycleError struct{}

func (c *cycleError) Error() string { return "cycle" }
func (c *cycleError) Unwrap() error { return c }

type cycleJoinError struct{}

func (c *cycleJoinError) Error() string   { return "cycle" }
func (c *cycleJoinError) Unwrap() []error { return []error{c} }

type cycleAError struct{ b *cycleBError }

func (c *cycleAError) Error() string { return "a" }
func (c *cycleAError) Unwrap() error { return c.b }

type cycleBError struct{ a *cycleAError }

func (c *cycleBError) Error() string { return "b" }
func (c *cycleBError) Unwrap() error { return c.a }

func TestAllCycle(t *testing.T) {
	count := 0
	for range All(&cycleError{}) {
		count++
	}
	assert.Equal(t, count, 1)
}

func TestAllCycleJoin(t *testing.T) {
	count := 0
	for range All(&cycleJoinError{}) {
		count++
	}
	assert.Equal(t, count, 1)
}

func TestAllCyclePair(t *testing.T) {
	a := &cycleAError{b: &cycleBError{}}
	a.b.a = a
	count := 0
	for range All(a) {
		count++
	}
	assert.Equal(t, count, 2)
}

func TestAllAllocs(t *testing.T) {
	err := newTestError()
	assert.AllocsPerRun(t, 100, func() {
		for range All(err) {
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
	m := FirstKeys(seq)
	assert.MapEqual(t, m, map[string]int{
		"a": 1,
		"b": 2,
		"c": 4,
	})
}

func TestFirstKeysEmpty(t *testing.T) {
	seq := func(yield func(string, int) bool) {}
	m := FirstKeys(seq)
	assert.MapNil(t, m)
}

func BenchmarkAll(b *testing.B) {
	err := newTestError()
	for b.Loop() {
		for range All(err) {
		}
	}
}
