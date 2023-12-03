package xerrors

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

type joinError struct {
	errs []error
}

// LogValue implements the [slog.LogValuer] interface
// it is our main point of entry to format the error as an attribute of a [slog.Record].
func (err *joinError) LogValue() slog.Value {
	info := Info(err)

	return slog.GroupValue(
		slog.String("message", info.ErrorChain),
		slog.String("stacktrace", strings.Join(info.StackTraces, "\n")),
		slog.Any("values", info.Values),
	)
}

// Format implements the [fmt.Formatter] interface
// it is our main point of entry to format the error using the [fmt] package.
func (err *joinError) Format(s fmt.State, verb rune) {
	format(err, s, verb)
}

func (err *joinError) message(verbose bool) string {
	if !verbose {
		return err.fmt("%s", "%s", ": ", " + ")
	}

	return err.fmt("%+v", "%+v", ": ", "\n")
}

func (err *joinError) fmt(mainErrFmt, secErrFmt, mainErrSep, secErrSep string) string {
	mainErrStr := fmt.Sprintf(mainErrFmt, err.errs[0])
	builder := bufferPool.Get().(*strings.Builder)
	builder.Reset()
	builder.WriteString(mainErrStr)
	var hasErrors bool
	errors := err.errs[1:]
	for i := range errors {
		s := fmt.Sprintf(secErrFmt, errors[i])
		if s != "" {
			if hasErrors {
				builder.WriteString(secErrSep)
			} else {
				builder.WriteString(mainErrSep)
			}
			builder.WriteString(s)
			hasErrors = true
		}
	}

	result := builder.String()
	bufferPool.Put(builder)

	return result
}

// stringify implement error.
func (err *joinError) Error() string {
	return stringify(err)
}

func (err *joinError) Unwrap() []error {
	n := 1 + len(err.errs) // Length of the resulting slice
	i := make([]error, 0, n)

	i = append(i, err.errs...)

	return i
}

// JoinStack creates a new error that represents an error chain by joining the original error
// with a list of additional errors.
// a stack is added if one is not already present in the error chain.
func JoinStack(errs ...any) error {
	n := 0
	for _, err := range errs {
		switch err.(type) {
		case string:
			n++
		case error:
			n++
		}
	}

	e := &joinError{
		errs: make([]error, 0, n),
	}
	var stacked bool
	for i := 0; i < len(errs); i++ {
		switch s := errs[i].(type) {
		case string:
			e.errs = append(e.errs, newErrorString(s))
		case error:
			if s != nil {
				e.errs = append(e.errs, s)
				b := hasStack(s)
				if b {
					stacked = true
				}
			}
		}
	}
	if !stacked {
		e.errs = append(e.errs, WithStackSkip(2))
	}
	return e
}

// Join creates a new error that represents an error chain by joining the original error
// with a list of additional errors.
func Join(errs ...any) error {
	n := 0
	for _, err := range errs {
		switch err.(type) {
		case string:
			n++
		case error:
			n++
		}
	}

	e := &joinError{
		errs: make([]error, 0, n),
	}

	for i := 0; i < len(errs); i++ {
		switch s := errs[i].(type) {
		case string:
			e.errs = append(e.errs, newErrorString(s))
		case error:
			if s != nil {
				e.errs = append(e.errs, s)
			}
		}
	}
	return e
}
