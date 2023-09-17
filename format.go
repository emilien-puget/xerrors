package xerrors

import (
	"fmt"
	"io"
	"sync"
)

// bytePool ptr on []bytes because of SA6002
// A sync.Pool is used to avoid unnecessary allocations and reduce the amount of work the garbage collector has to do.
//
// When passing a value that is not a pointer to a function that accepts an interface, the value needs to be placed on the heap, which means an additional allocation.
// Slices are a common thing to put in sync.Pools, and they’re structs with 3 fields (length, capacity, and a pointer to an array).
// In order to avoid the extra allocation, one should store a pointer to the slice instead.
var bytePool = sync.Pool{
	New: func() interface{} {
		bytes := make([]byte, 0, 1024)
		return &bytes
	},
}

type formattable interface {
	error
	message(verbose bool) string
}

func stringify(err formattable) string {
	message := err.message(false)
	werr := Unwrap(err)

	if message != "" && werr != nil {
		res := bytePool.Get().(*[]byte)
		result := *res
		result = result[:0]
		result = append(result, message...)
		result = append(result, ": "...)
		result = append(result, werr.Error()...)
		s := string(result)
		bytePool.Put(&result)
		return s
	}

	if message != "" {
		return message
	}

	if werr != nil {
		res := bytePool.Get().(*[]byte)
		result := *res
		result = result[:0]
		result = result[:0]
		result = append(result, werr.Error()...)
		s := string(result)
		bytePool.Put(&result)
		return s
	}

	return ""
}

func format(err formattable, s fmt.State, verb rune) {
	switch {
	case verb == 'v' && s.Flag('+'):
		werr := Unwrap(err)
		if werr != nil {
			_, _ = fmt.Fprintf(s, "%+v", werr)
		}
		msg := err.message(true)
		if msg != "" {
			_, _ = fmt.Fprintf(s, "%s", msg)
		}
	case verb == 'v' || verb == 's':
		_, _ = io.WriteString(s, err.Error())
	case verb == 'q':
		_, _ = fmt.Fprintf(s, "%q", err.Error())
	}
}
