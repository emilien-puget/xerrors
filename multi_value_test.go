package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiValue(t *testing.T) {
	err := New("error")
	err = Join(err, WithValues(map[string]any{"foo": "bar"}))
	err = Join(err, "wrapped")
	vals := Values(err)
	expected := map[string]interface{}{
		"foo": "bar",
	}
	assert.Equal(t, expected, vals)
}

func TestMultiValueOverWrite(t *testing.T) {
	err := New("error")
	err = Join(err, WithValues(map[string]any{"test": 1}))
	err = Join(err, WithValues(map[string]any{"test": 2}))
	err = Join(err, "wrapped")
	vals := Values(err)
	expected := map[string]interface{}{
		"test": 1,
	}
	assert.Equal(t, expected, vals)
}
