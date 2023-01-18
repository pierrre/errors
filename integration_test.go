package errors_test

import (
	"regexp"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errignore"
	"github.com/pierrre/errors/errtag"
	"github.com/pierrre/errors/errtmp"
	"github.com/pierrre/errors/errval"
	"github.com/pierrre/errors/errverbose"
)

func TestIntegration(t *testing.T) {
	err := errors.New("error")
	err = errors.Wrap(err, "test")
	err = errignore.Wrap(err)
	err = errtmp.Wrap(err, true)
	err = errtag.Wrap(err, "a", "b")
	err = errval.Wrap(err, "c", "d")
	t.Run("Error", func(t *testing.T) {
		s := err.Error()
		expected := "test: error"
		if s != expected {
			t.Fatalf("unexpected message: got %q, want %q", s, expected)
		}
	})
	t.Run("Verbose", func(t *testing.T) {
		s := errverbose.String(err)
		expected := regexp.MustCompile(`^test: error\nvalue c = d\ntag a = b\ntemporary = true\nignored\nstack\n(\t.+ .+:\d+\n)+\n$`)
		if !expected.MatchString(s) {
			t.Fatalf("unexpected verbose message:\ngot: %q\nwant match: %q", s, expected)
		}
	})
	t.Run("Stack", func(t *testing.T) {
		fs := errors.StackFrames(err)
		if len(fs) != 1 {
			t.Fatalf("unexpected stack frames length: got %d, want 1", len(fs))
		}
	})
	t.Run("Ignore", func(t *testing.T) {
		if !errignore.Is(err) {
			t.Fatal("not ignored")
		}
	})
	t.Run("Temporary", func(t *testing.T) {
		if !errtmp.Is(err) {
			t.Fatal("not temporary")
		}
	})
	t.Run("Tag", func(t *testing.T) {
		tags := errtag.Get(err)
		if len(tags) != 1 {
			t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
		}
		if tags["a"] != "b" {
			t.Fatalf("unexpected tag: got %q, want %q", tags["a"], "b")
		}
	})
	t.Run("Value", func(t *testing.T) {
		values := errval.Get(err)
		if len(values) != 1 {
			t.Fatalf("unexpected length: got %d, want %d", len(values), 1)
		}
		if values["c"] != "d" {
			t.Fatalf("unexpected tag: got %v, want %q", values["c"], "d")
		}
	})
}
