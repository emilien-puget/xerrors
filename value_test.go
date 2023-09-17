package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	err := New("error")
	err = Join(err, WithValue("foo", "bar"))
	err = Join(err, "wrapped")
	vals := Values(err)
	expected := map[string]interface{}{
		"foo": "bar",
	}
	assert.Equal(t, expected, vals)
}

func TestValueOverWrite(t *testing.T) {
	err := New("error")
	err = Join(err, WithValue("test", 1), WithValue("test", 2))
	err = Join(err, "wrapped")
	vals := Values(err)
	expected := map[string]interface{}{
		"test": 1,
	}
	assert.Equal(t, expected, vals)
}
