package xerrors

// Values returns the values associated to an error.
func Values(err error) map[string]any {
	errors := FlattenErrors(err)
	vals := make(map[string]any)
	for i := range errors {
		switch et := errors[i].(type) {
		case Valuer:
			key, v := et.Value()
			if _, ok := vals[key]; ok {
				continue
			}
			vals[key] = v
		case MultiValuer:
			for s, a := range et.Value() {
				if _, ok := vals[s]; ok {
					continue
				}
				vals[s] = a
			}
		}
	}
	return vals
}
