package strconvio

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/pierrre/assert"
)

func TestWriteInt(t *testing.T) {
	for _, tc := range []struct {
		i        int64
		expected string
	}{
		{
			i:        0,
			expected: "0",
		},
		{
			i:        1,
			expected: "1",
		},
		{
			i:        2,
			expected: "2",
		},

		{
			i:        1234567890,
			expected: "1234567890",
		},
		{
			i:        -1,
			expected: "-1",
		},
		{
			i:        -1234567890,
			expected: "-1234567890",
		},
	} {
		t.Run(strconv.FormatInt(tc.i, 10), func(t *testing.T) {
			buf := new(bytes.Buffer)
			n, err := WriteInt(buf, tc.i, 10)
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expected), n)
			assert.Equal(t, tc.expected, buf.String())
		})
	}
}

func BenchmarkWriteInt(b *testing.B) {
	for _, i := range []int64{
		0,
		1,
		2,
		1234567890,
		-1,
		-1234567890,
	} {
		b.Run(strconv.FormatInt(i, 10), func(b *testing.B) {
			buf := new(bytes.Buffer)
			for n := 0; n < b.N; n++ {
				_, _ = WriteInt(buf, i, 10)
				buf.Reset()
			}
		})
	}
}
