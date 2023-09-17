package xerrors

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBase(t *testing.T) {
	err := New("error")
	sprint := fmt.Sprint(err)
	assert.Equal(t, "error", sprint)
}

func TestNewUnwrap(t *testing.T) {
	err := New("error")
	assert.ErrorIs(t, err, &stack{})
	newErr := Unwrap(err)
	assert.NotErrorIs(t, newErr, &stack{})
}

func TestIs(t *testing.T) {
	err := io.EOF
	err = Join(err, "test")
	assert.True(t, Is(err, io.EOF))
}

func TestAs(t *testing.T) {
	err := &net.ParseError{
		Type: "_type",
		Text: "_text",
	}
	wrapped := Join(err, "test")

	var expected *net.ParseError
	assert.True(t, As(wrapped, &expected))
}
