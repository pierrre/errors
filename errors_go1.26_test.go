//go:build go1.26

package errors_test

import (
	"io/fs"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
)

func TestAsType(t *testing.T) {
	err := errbase.New("error")
	err = &fs.PathError{Err: err}
	err = Wrap(err, "test")
	pathError, ok := AsType[*fs.PathError](err)
	assert.True(t, ok)
	assert.NotZero(t, pathError)
}
