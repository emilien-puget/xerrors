package xerrors

import "fmt"

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
	return fmt.Sprintf("value: %s %q", err.key, err.value)
}

// Values returns the values associated to an error.
func Values(err error) map[string]any {
	errors := FlattenErrors(err)
	vals := make(map[string]any)
	for i := range errors {
		err, ok := errors[i].(*value)
		if !ok {
			continue
		}
		k, v := err.Value()
		_, ok = vals[k]
		if ok {
			continue
		}
		vals[k] = v
	}
	return vals
}
