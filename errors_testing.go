//go:build !errorsnotesting

package errors

import (
	"testing"
)

func init() {
	var f func(error)
	if testing.Testing() {
		f = func(err error) {
			panic(err)
		}
	}
	ReportGlobalInit.Store(f)
}
