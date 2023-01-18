# Errors

Go errors library.

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/errors.svg)](https://pkg.go.dev/github.com/pierrre/errors)

## Features

- [Stack trace](#stack-trace)
- [Message](#message)
- [Verbose message](#verbose-message)
- [Drop-in replacement of the std `errors` package](#migrate-from-the-std-errors-package)
- [Easy to extend](#extend)

## Message

[`New()`](https://pkg.go.dev/github.com/pierrre/errors#New) creates an error with a message and a stack.

```go
err := errors.New("error")
```

[`Wrap()`](https://pkg.go.dev/github.com/pierrre/errors#Wrap) adds a message to an error, and optionally adds a stack if the error doesn't have one.

```go
err = errors.Wrap(err, "message")
```

## Stack trace

[`Stack()`](https://pkg.go.dev/github.com/pierrre/errors#Stack) adds a stack to an error. This is only useful if the wrapped error has a stack from a different goroutine. Most applications should use [`Wrap()`](https://pkg.go.dev/github.com/pierrre/errors#Wrap) instead.

```go
err = errors.Stack(err)
```

[`StackFrames()`](https://pkg.go.dev/github.com/pierrre/errors#StackFrames) returns the [stack frames](https://pkg.go.dev/runtime#Frames) of the error.

```go
frames := errors.StackFrames(err)
```

## Verbose message

The error verbose message shows additional information about the error.
Wrapping functions may provide a verbose message (stack, tag, value, etc.)

The [`errverbose`](https://pkg.go.dev/github.com/pierrre/errors/errverbose) package provides utilities to manage error verbose messages. The [`Write()`](https://pkg.go.dev/github.com/pierrre/errors/errverbose#Write)/[`String()`](https://pkg.go.dev/github.com/pierrre/errors/errverbose#String)/[`Formatter()`](https://pkg.go.dev/github.com/pierrre/errors/errverbose#Formatter) functions write/return/format the error verbose message.

The first line is the error's message.
The following lines are the verbose message of the error chain.

Example:

```text
test: error
value c = d
tag a = b
temporary = true
ignored
stack
    github.com/pierrre/errors_test.TestIntegration integration_test.go:15
    testing.tRunner testing.go:1446
    runtime.goexit asm_amd64.s:1594
```

## Extend

Create a custom error type:

- Create a type implementing the [`error`](https://pkg.go.dev/builtin#error) interface
- Optionally implement the [`Unwrap() error`](https://pkg.go.dev/errors#Unwrap) method
- Optionally implement the [`errverbose.Interface`](https://pkg.go.dev/github.com/pierrre/errors/errverbose#Interface) interface

See the provided packages as example:

- [`errbase`](https://pkg.go.dev/github.com/pierrre/errors/errbase): create a base error (e.g. sentinel error)
- [`errmsg`](https://pkg.go.dev/github.com/pierrre/errors/errmsg): add a message to an error
- [`errtag`](https://pkg.go.dev/github.com/pierrre/errors/errtag): add a tag to an error
- [`errval`](https://pkg.go.dev/github.com/pierrre/errors/errval): add a value to an error
- [`errignore`](https://pkg.go.dev/github.com/pierrre/errors/errignore): mark an error as ignored
- [`errtmp`](https://pkg.go.dev/github.com/pierrre/errors/errtmp): mark an error as temporary

## Migrate from the std `errors` package

- Replace the import `errors` with `github.com/pierrre/errors`
- Replace `fmt.Errorf("some wessage: %w", err)` with `errors.Wrap(err, "some message")`
