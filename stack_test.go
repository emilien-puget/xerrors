package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	err := New("error")
	sfs := StackFrames(err)

	frames := sfs.Frames()
	assert.Len(t, frames, 3)
	assert.Equal(t, frames[0].Function, "github.com/emilien-puget/xerrors.TestStack")
	assert.Equal(t, frames[1].Function, "testing.tRunner")
	assert.Equal(t, frames[2].Function, "runtime.goexit")
}

func TestStackNil(t *testing.T) {
	err := withStack(nil, 0)
	assert.NoError(t, err)
}
