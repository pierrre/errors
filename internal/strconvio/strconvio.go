// Package strconvio provides utilities to write values string representation to a writer.
package strconvio

import (
	"io"
	"strconv"
	"sync"
)

var bytesPool = sync.Pool{
	New: func() any {
		var v []byte
		return &v
	},
}

// WriteInt writes the string representation of the integer to the writer.
func WriteInt(w io.Writer, i int64, base int) (int, error) {
	if 0 <= i && i < 100 && base == 10 {
		return io.WriteString(w, strconv.FormatInt(i, base)) //nolint:wrapcheck // It's fine.
	}
	bp := bytesPool.Get().(*[]byte) //nolint:forcetypeassert // The pool only contains *[]byte.
	defer bytesPool.Put(bp)
	*bp = strconv.AppendInt((*bp)[:0], i, base)
	return w.Write(*bp) //nolint:wrapcheck // It's fine.
}
