package xerrors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type StackFramess []uintptr

func (s StackFramess) Frames() []Frame {
	r := make([]Frame, len(s))
	f := runtime.CallersFrames(s)
	n := 0

	for more := true; more; {
		var frame runtime.Frame
		frame, more = f.Next()
		r[n] = Frame{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		}
		n++
	}
	return r
}

type Frame struct {
	File     string
	Line     int
	Function string
}

func (s Frame) String() string {
	builder := bufferPool.Get().(*strings.Builder)
	builder.Reset()
	defer bufferPool.Put(builder)

	_, _ = fmt.Fprintf(builder, "%s %s:%d", s.Function, s.File, s.Line)
	return builder.String()
}

type stack struct {
	err     error
	callers StackFramess
}

func (err *stack) Error() string {
	return stringify(err)
}

func (err *stack) message(verbose bool) string {
	if !verbose {
		return ""
	}

	builder := bufferPool.Get().(*strings.Builder)
	builder.Reset()
	defer bufferPool.Put(builder)

	builder.WriteString("stack\n")
	framesString := StackFrames(err)
	for _, fr := range framesString.Frames() {
		builder.WriteString("\t" + fr.String() + "\n")
	}

	return "\n" + builder.String()
}

// Is implements the anonymous interface Is
// it provides an alternative to As that allocates less memory if you are not interested in the actual stacktrace.
func (err *stack) Is(target error) bool {
	_, ok := target.(*stack)
	return ok
}

func (err *stack) Unwrap() error {
	return err.err
}

func (err *stack) StackFrames() StackFramess {
	return err.callers
}

// Format implements the [fmt.Formatter] interface
// it is our main point of entry to format the error using the [fmt] package.
func (err *stack) Format(s fmt.State, verb rune) {
	format(err, s, verb)
}

// StackFrames returns the list of *StackFramess associated to an error.
func StackFrames(err error) *StackFramess {
	var fss *StackFramess
	for ; err != nil; err = errors.Unwrap(err) {
		var errS *stack
		ok := errors.As(err, &errS)
		if ok {
			fs := errS.StackFrames()
			return &fs
		}

	}
	return fss
}

func callers(skip int) StackFramess {
	pc := make([]uintptr, 32)
	n := runtime.Callers(skip+1, pc)
	return pc[:n]
}

func ensureStack(err error, skip int) error {
	if !hasStack(err) {
		err = withStack(err, skip+1)
	}
	return err
}

func hasStack(err error) bool {
	return Is(err, &stack{})
}

func withStack(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &stack{
		err:     err,
		callers: callers(skip + 1),
	}
}
