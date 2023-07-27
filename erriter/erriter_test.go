package erriter_test

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/errors/errmsg"
)

func TestIter(t *testing.T) {
	err := errbase.New("error")
	err = errors.Join(err, err)
	err = errmsg.Wrap(err, "test")
	count := 0
	erriter.Iter(err, func(err error) {
		count++
		assert.Error(t, err)
	})
	assert.Equal(t, count, 4)
}
