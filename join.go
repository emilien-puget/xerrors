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
	err  error
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

	return err.fmt("%+v", "\t%+v", ": ", "\n")
}

func (err *joinError) fmt(mainErrFmt, secErrFmt, mainErrSep, secErrSep string) string {
	mainErrStr := fmt.Sprintf(mainErrFmt, err.err)

	builder := bufferPool.Get().(*strings.Builder)
	builder.Reset()
	builder.WriteString(mainErrStr)

	var hasErrors bool
	for i := range err.errs {
		s := fmt.Sprintf(secErrFmt, err.errs[i])
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

	i = append(i, err.err)
	i = append(i, err.errs...)

	return i
}

// Join creates a new error that represents an error chain by joining the original error
// with a list of additional errors.
func Join(ogErr error, errs ...any) error {
	n := 0
	for _, err := range errs {
		switch err.(type) {
		case string:
			n++
		case error:
			n++
		}
	}

	stackedErr := ensureStack(ogErr, 2)
	e := &joinError{
		err:  stackedErr,
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
