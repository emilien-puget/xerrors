package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorF(t *testing.T) {
	errErrorF := New("errorF")
	newErr := ErrorF("wrapping: %w", errErrorF)

	assert.ErrorIs(t, newErr, errErrorF)

	var joinErr *joinError
	require.ErrorAs(t, newErr, &joinErr)

	assert.ErrorIs(t, newErr, errErrorF)
	assert.Len(t, joinErr.errs, 2)
	var stacked *stack
	require.ErrorAs(t, newErr, &stacked)
}
