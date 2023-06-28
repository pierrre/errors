package pierrrepretty

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errval"
	"github.com/pierrre/errors/errverbose"
)

func init() {
	ConfigureDefault()
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = errval.Wrap(err, "foo", "bar")
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	s := v.ErrorVerbose()
	assert.Equal(t, s, "value foo = (string) (len=3) \"bar\"")
}
