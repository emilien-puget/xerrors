package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiValue(t *testing.T) {
	err := New("error")
	err = JoinStack(err, WithValues(map[string]any{"foo": "bar"}))
	err = JoinStack(err, "wrapped")
	vals := Values(err)
	expected := map[string]interface{}{
		"foo": "bar",
	}
	assert.Equal(t, expected, vals)
}

func TestMultiValueOverWrite(t *testing.T) {
	err := New("error")
	err = JoinStack(err, WithValues(map[string]any{"test": 1}))
	err = JoinStack(err, WithValues(map[string]any{"test": 2}))
	err = JoinStack(err, "wrapped")
	vals := Values(err)
	expected := map[string]interface{}{
		"test": 1,
	}
	assert.Equal(t, expected, vals)
}
