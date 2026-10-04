package erriter_test

import (
	"slices"
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

type multiError struct{ errs []error }

func (m *multiError) Error() string   { return "multi" }
func (m *multiError) Unwrap() []error { return m.errs }

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

// newCyclePair returns two errors that reference each other.
func newCyclePair() (a *cycleAError, b *cycleBError) {
	b = &cycleBError{}
	a = &cycleAError{b: b}
	a.b.a = a
	return
}

// collect returns the sequence of errors yielded by All(err).
func collect(err error) []error {
	return slices.Collect(All(err))
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

func TestAllNil(t *testing.T) {
	count := 0
	for range All(nil) {
		count++
	}
	assert.Equal(t, count, 0)
}

func TestAllNilChild(t *testing.T) {
	// A multi error may unwrap to a list containing nil.
	leaf := errbase.New("leaf")
	root := &multiError{errs: []error{nil, leaf}}
	got := collect(root)
	assert.Equal(t, len(got), 2)
	assert.Equal(t, got[0], error(root))
	assert.Equal(t, got[1], leaf)
}

func TestAllOrder(t *testing.T) {
	a := errbase.New("a")
	b := errbase.New("b")
	c := errbase.New("c")
	wa := errors.Wrap(a, "wa")
	root := errors.Join(wa, b, c)
	got := collect(root)
	// The sequence is: root, its join, then for each child its whole
	// subtree: the wrapped chain (wa, its stack, a), then b, then c.
	assert.Equal(t, len(got), 7)
	assert.Equal(t, got[0], root)
	assert.Equal(t, got[2], wa)
	assert.Equal(t, got[4], a)
	assert.Equal(t, got[5], b)
	assert.Equal(t, got[6], c)
}

func TestAllDiamond(t *testing.T) {
	// A leaf reachable through two different paths is visited once per path.
	e := errbase.New("e")
	got := collect(errors.Join(e, e))
	assert.Equal(t, len(got), 4)
	assert.Equal(t, got[2], e)
	assert.Equal(t, got[3], e)
}

func TestAllDeepChain(t *testing.T) {
	// A long chain must not overflow the stack.
	const n = 100000
	err := errbase.New("leaf")
	for range n {
		err = errmsg.Wrap(err, "w")
	}
	count := 0
	for range All(err) {
		count++
	}
	assert.Equal(t, count, n+1)
}

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
	a, _ := newCyclePair()
	count := 0
	for range All(a) {
		count++
	}
	assert.Equal(t, count, 2)
}

func TestAllCycleSibling(t *testing.T) {
	// A cycle is cut, but iteration continues with the next sibling.
	a, b := newCyclePair()
	leaf := errbase.New("leaf")
	root := errors.Join(a, leaf)
	got := collect(root)
	assert.Equal(t, len(got), 5)
	assert.Equal(t, got[2], error(a))
	assert.Equal(t, got[3], error(b))
	assert.Equal(t, got[4], leaf)
}

func TestAllAllocs(t *testing.T) {
	err := newTestError()
	assert.AllocsPerRun(t, 100, func() {
		for range All(err) {
		}
	}, 0)
}

func BenchmarkAll(b *testing.B) {
	err := newTestError()
	for b.Loop() {
		for range All(err) {
		}
	}
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
