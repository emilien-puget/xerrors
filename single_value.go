package xerrors

import "fmt"

// Valuer is an interface that allows custom error types to associate a single key-value pair with an error.
// Implement this interface in your custom error type to provide specific metadata for the error.
// If a key already exists within the error, adding a new value for that key will not replace the existing value.
type Valuer interface {
	// Value returns the key and value associated with the error.
	Value() (key string, value any)
}

// WithValue return an error that contains a value.
func WithValue(key string, val any) error {
	return &value{
		key:   key,
		value: val,
	}
}

type value struct {
	key   string
	value any
}

func (err *value) Value() (key string, value any) {
	return err.key, err.value
}

// Format implements the [fmt.Formatter] interface
// it is our main point of entry to format the error using the [fmt] package.
func (err *value) Format(s fmt.State, verb rune) {
	format(err, s, verb)
}

func (err *value) Error() string {
	return stringify(err)
}

func (err *value) message(verbose bool) string {
	if !verbose {
		return ""
	}
	return fmt.Sprintf("value: %s \"%v\"", err.key, err.value)
}
