package xerrors

import (
	"fmt"
	"strings"
)

// MultiValuer is an interface that enables custom error types to associate multiple key-value pairs with an error.
// Implement this interface in your custom error type when you need to attach various context information to the error.
// If a key already exists within the error, adding a new value for that key will not replace the existing value.
type MultiValuer interface {
	// Value returns a map of key-value pairs associated with the error.
	Value() map[string]any
}

// WithValues return an error that contains multiple values.
func WithValues(v map[string]any) error {
	return &multiValue{
		values: v,
	}
}

type multiValue struct {
	values map[string]any
}

func (err *multiValue) Value() map[string]any {
	return err.values
}

// Format implements the [fmt.Formatter] interface
// it is our main point of entry to format the error using the [fmt] package.
func (err *multiValue) Format(s fmt.State, verb rune) {
	format(err, s, verb)
}

func (err *multiValue) Error() string {
	return stringify(err)
}

func (err *multiValue) message(verbose bool) string {
	if !verbose {
		return ""
	}
	var keyValuePairs []string

	for key, v := range err.values {
		// Format the key-value pair and add it to the slice.
		keyValuePairs = append(keyValuePairs, fmt.Sprintf("%s: \"%v\"", key, v))
	}

	// Join the key-value pairs into a single string, separated by spaces.
	return "values: [" + strings.Join(keyValuePairs, " ") + "]"
}
