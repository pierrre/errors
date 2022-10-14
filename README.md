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

[`New`](https://pkg.go.dev/github.com/pierrre/errors#New) and [`Newf`](https://pkg.go.dev/github.com/pierrre/errors#Newf) functions create an error with a message and a stack.

```go
err := errors.New("error")
```

```go
err := errors.Newf("error %s", "foo")
```

[`Wrap`](https://pkg.go.dev/github.com/pierrre/errors#Wrap) and [`Wrapf`](https://pkg.go.dev/github.com/pierrre/errors#Wrapf) functions add a message to an error, and optionally add a stack if the error doesn't have one.

```go
err = errors.Wrap(err, "message")
```

```go
err = errors.Wrapf(err, "message %d", 1)
```

[`Message`](https://pkg.go.dev/github.com/pierrre/errors#Message) and [`Messagef`](https://pkg.go.dev/github.com/pierrre/errors#Messagef) functions add a message to an error. Most applications should use [`Wrap`](https://pkg.go.dev/github.com/pierrre/errors#Wrap) and [`Wrapf`](https://pkg.go.dev/github.com/pierrre/errors#Wrapf) instead, because they automatically add a stack.

```go
err = errors.WithMessage(err, "message")
```

```go
err = errors.WithMessagef(err, "message %d", 1)
```

## Stack trace

[`Stack`](https://pkg.go.dev/github.com/pierrre/errors#Stack) function adds a stack to an error. This is only useful if the wrapped error has a stack from a different goroutine. Most applications should use [`Wrap`](https://pkg.go.dev/github.com/pierrre/errors#Wrap) and [`Wrapf`](https://pkg.go.dev/github.com/pierrre/errors#Wrapf) instead.

```go
err = errors.Stack(err)
```

[`StackFrames`](https://pkg.go.dev/github.com/pierrre/errors#StackFrames) function returns the [stack frames](https://pkg.go.dev/runtime#Frames) of the error.

```go
frames := errors.StackFrames(err)
```

## Verbose message

The error verbose message shows additional information about the error.
Wrapping functions may provide a verbose message (stack, tag, value, etc.)

The [`Verbose`](https://pkg.go.dev/github.com/pierrre/errors#Verbose)/[`VerboseString`](https://pkg.go.dev/github.com/pierrre/errors#VerboseString)/[`VerboseFormatter`](https://pkg.go.dev/github.com/pierrre/errors#VerboseFormatter) functions write/return/format the error verbose message.

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
- Optionally implement the [`Verboser`](https://pkg.go.dev/github.com/pierrre/errors#Verboser) interface

See the provided packages as example:

- [`errtag`](https://pkg.go.dev/github.com/pierrre/errors/errtag): add a tag to an error
- [`errval`](https://pkg.go.dev/github.com/pierrre/errors/errval): add a value to an error
- [`errignore`](https://pkg.go.dev/github.com/pierrre/errors/errignore): mark an error as ignored
- [`errtmp`](https://pkg.go.dev/github.com/pierrre/errors/errtmp): mark an error as temporary

## Migrate from the std `errors` package

- Replace the import `errors` with `github.com/pierrre/errors`
- Replace `fmt.Errorf("some wessage: %w", err)` with `errors.Wrap(err, "some message")`
