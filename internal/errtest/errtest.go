// Package errtest provides internal tools for tests.
package errtest

import (
	"github.com/pierrre/assert/ext/pierrrecompare"
	"github.com/pierrre/assert/ext/pierrreerrors"
	"github.com/pierrre/assert/ext/pierrrepretty"
)

// Configure configures tools used in tests.
func Configure() {
	pierrrecompare.Configure()
	pierrrepretty.ConfigureDefault()
	pierrreerrors.Configure()
}
