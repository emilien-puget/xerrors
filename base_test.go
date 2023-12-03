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

func TestIs(t *testing.T) {
	err := io.EOF
	err = JoinStack(err, "test")
	assert.True(t, Is(err, io.EOF))
}

func TestAs(t *testing.T) {
	err := &net.ParseError{
		Type: "_type",
		Text: "_text",
	}
	wrapped := JoinStack(err, "test")

	var expected *net.ParseError
	assert.True(t, As(wrapped, &expected))
	assert.ErrorIs(t, wrapped, &stack{})
}
