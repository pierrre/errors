package errbase_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func Test(t *testing.T) {
	err := Newf("error %d", 1)
	assert.ErrorEqual(t, err, "error 1")
}

func Example() {
	err := New("error")
	s := err.Error()
	fmt.Println(s)
	// Output: error
}
