# Xerrors

[![codecov](https://codecov.io/gh/emilien-puget/xerrors/graph/badge.svg?token=kjMuQevvTc)](https://codecov.io/gh/emilien-puget/xerrors)
[![Test and coverage](https://github.com/emilien-puget/xerrors/actions/workflows/codecov.yml/badge.svg)](https://github.com/emilien-puget/xerrors/actions/workflows/codecov.yml)

## Introduction

The `xerrors` package aims to enhance error handling in Go programs by providing custom error types and utility
functions for error creation and manipulation. It includes features such as:

- Custom error types
- Error chaining
- Error value association
- Stack trace capture
- Loggable

## Usage

```go
import "github.com/emilien-puget/xerrors"
```

### Creating New Errors

You can create a new error with a message using the `New` function:

```go
err := xerrors.New("This is an error message")
```

### Chaining Errors

You can chain multiple errors together using the `Join` function:

```go
err1 := xerrors.New("Error 1")
err2 := xerrors.New("Error 2")
chainedErr := xerrors.Join(err1, err2)
```

### Logging Errors

Errors created using this package implement the slog.Valuer interface. When such an error is logged using slog from the
standard library, it will automatically unpack the error with all available information, including stack traces and
associated values.

### Getting Error Information

To retrieve information about an error, such as its message, stack traces, and associated values, you can use the `Info`
function:

```go
err := xerrors.New("An error occurred")
info := xerrors.Info(err)
fmt.Println("Error Message:", info.ErrorChain)
fmt.Println("Stack Traces:", info.StackTraces)
fmt.Println("Values:", info.Values)
```

### Associating Values with Errors

You can associate values with errors using the `WithValue` function:

```go
valueErr := xerrors.WithValue("key", "value")
chainedErrWithValue := xerrors.Join(chainedErr, valueErr)
```

The xerrors package also exposes two interfaces, ``Valuer`` and ``MultiValuer``, that allow you to create your own
custom error types that can carry associated values.

```go
package main

type Valuer interface {
	Value() (key string, value any)
}

type MultiValuer interface {
	Value() map[string]any
}
```

#### Adding Values for Existing Keys

When using the `Valuer` and `MultiValuer` interfaces to associate values with errors, it's important to note that if a
key already exists within the error, adding a new value for that key will not replace the existing value.
Instead, the new value will be associated with the existing key, allowing you to accumulate context information without
overwriting
valuable data already associated with the error.

This behavior ensures that each piece of metadata you attach to an error is preserved.

### Checking Error Relationships

To check if one error is related to another, you can use functions like Is and As. Please note that Is and As methods
provided by xerrors are purely aliases to the standard library, and their purpose is to ensure that you use this
package's error handling mechanisms instead of the standard library's:

```go
err1 := xerrors.New("Error 1")
err2 := xerrors.New("Error 2")
isRelated := xerrors.Is(err1, err2)
```

aliases to

## Formatting Errors

Errors can be formatted using the `fmt` package. By default, errors are formatted as strings. To include additional
information, you can use the `+` verb for a more detailed representation.

```go
err := xerrors.New("An error occurred")
fmt.Printf("%+v\n", err) // Detailed error information
fmt.Printf("%s\n", err) // Basic error message
```

## Stack Traces

The package captures stack traces for errors. Stack traces can be retrieved using the `StackFrames` function and
formatted as strings.
