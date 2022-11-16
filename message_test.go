package errors

import (
	"fmt"
	"testing"
)

func TestMessage(t *testing.T) {
	err := New("error")
	err = Message(err, "test")
	s := err.Error()
	expected := "test: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestMessageNil(t *testing.T) {
	err := Message(nil, "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessageEmpty(t *testing.T) {
	err := New("error")
	err = Message(err, "")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func ExampleMessage() {
	err := New("error")
	err = Message(err, "message")
	fmt.Println(err)
	// Output: message: error
}
