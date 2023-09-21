package xerrors

type ErrorInfo struct {
	ErrorChain  string
	StackTraces []string
	Values      map[string]any
	Type        string
}

// Info returns information about the error chain, stack traces, and values.
func Info(err error) ErrorInfo {
	values := make(map[string]any)
	var stackTraces []string
	var errS *stack

	typeString := ""

	errors := FlattenErrors(err)
	for i := range errors {
		switch et := errors[i].(type) {
		case Valuer:
			key, v := et.Value()
			values[key] = v
		case MultiValuer:
			for s, a := range et.Value() {
				values[s] = a
			}
		case *stack:
			if errS == nil {
				errS = et
			}
		}
	}

	s := ""
	if fm, ok := err.(formattable); ok {
		s = stringify(fm)
	}

	if errS != nil {
		fs := errS.StackFrames()
		for _, frame := range fs.Frames() {
			stackTraces = append(stackTraces, frame.String())
		}
	}

	return ErrorInfo{
		ErrorChain:  s,
		StackTraces: stackTraces,
		Values:      values,
		Type:        typeString,
	}
}

func FlattenErrors(err error) []error {
	flatErrors := make([]error, 0)

	switch v := err.(type) {
	case interface{ Unwrap() error }:
		// The error implements Unwrap() error, unwrap it and add to flatErrors
		flatErrors = append(flatErrors, err)
		innerErr := v.Unwrap()
		if innerErr != nil {
			// Recursively call FlattenErrors for the inner error
			innerFlatErrors := FlattenErrors(innerErr)
			flatErrors = append(flatErrors, innerFlatErrors...)
		}
	case interface{ Unwrap() []error }:
		// The error implements Unwrap() []error, unwrap it and add its elements to flatErrors
		flatErrors = append(flatErrors, err)
		slice := v.Unwrap()
		for i := range slice {
			// Recursively call FlattenErrors for each inner error
			innerFlatErrors := FlattenErrors(slice[i])
			flatErrors = append(flatErrors, innerFlatErrors...)
		}
	default:
		// If the error doesn't implement either Unwrap interface, add it as is
		flatErrors = append(flatErrors, err)
	}

	return flatErrors
}
