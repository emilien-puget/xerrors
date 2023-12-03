package xerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack_JoinStack(t *testing.T) {
	err := JoinStack(New("error"))
	sfs := StackFrames(err)

	frames := sfs.Frames()
	assert.Len(t, frames, 3)
	assert.Equal(t, "github.com/emilien-puget/xerrors.TestStack_JoinStack", frames[0].Function)
	assert.Equal(t, "testing.tRunner", frames[1].Function)
	assert.Equal(t, "runtime.goexit", frames[2].Function)
}

func TestStack_Join(t *testing.T) {
	err := Join(New("error"), WithStack())
	sfs := StackFrames(err)

	frames := sfs.Frames()
	assert.Len(t, frames, 3)
	assert.Equal(t, "github.com/emilien-puget/xerrors.TestStack_Join", frames[0].Function)
	assert.Equal(t, "testing.tRunner", frames[1].Function)
	assert.Equal(t, "runtime.goexit", frames[2].Function)
}
