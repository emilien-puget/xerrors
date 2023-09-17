package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	err := New("plouf")
	err = Join(err, "its a wrap")
	err = Join(err, WithValue("key1", "value1"), WithValue("key2", "value2"), "An error occurred")

	info := Info(err)
	assert.Equal(t, "plouf: its a wrap: An error occurred", info.ErrorChain)
	assert.Equal(t, map[string]any{"key1": "value1", "key2": "value2"}, info.Values)
	assert.Len(t, info.StackTraces, 3)
}

func createErrorGraph(depth int) error {
	if depth <= 0 {
		return New("base error")
	}
	return Join(createErrorGraph(depth-1), "error message", WithValue("toto", 3))
}

func BenchmarkInfo(b *testing.B) {
	err := createErrorGraph(1000) // Adjust the depth as needed
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(err)
	}
}
