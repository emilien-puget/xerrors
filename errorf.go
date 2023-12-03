package xerrors

import basefmt "fmt"

func ErrorF(f string, err error) error {
	return basefmt.Errorf(f, EnsureStackSkip(err, 1))
}
