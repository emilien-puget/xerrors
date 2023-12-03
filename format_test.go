package xerrors

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	err := New("error")
	err = JoinStack(err, "its a wrap")
	err = JoinStack(err, "its another wrap", "with something more", WithValue("toto", "key1"), WithValues(map[string]any{"toto": "key1", "foo": 404}))
	for _, tc := range []struct {
		ft       string
		expected string
	}{
		{
			ft:       "%v",
			expected: regexp.QuoteMeta(`error: its a wrap: its another wrap + with something more`),
		},
		{
			ft:       "%+v",
			expected: "^error: its a wrap\nstack\n\tgithub.com/emilien-puget/xerrors.TestFormat ([^ ]+)format_test.go:13\n\ttesting.tRunner ([^ ]+)testing.go:1595\n\truntime.goexit ([^ ]+)asm_amd64.s:1650\n: its another wrap\nwith something more\nvalue: toto \"key1\"\nvalues: \\[toto: \"key1\" foo: \"404\"\\]$",
		},
		{
			ft:       "%s",
			expected: regexp.QuoteMeta("error: its a wrap: its another wrap + with something more"),
		},
		{
			ft:       "%q",
			expected: regexp.QuoteMeta("\"error: its a wrap: its another wrap + with something more\""),
		},
	} {
		t.Run(tc.ft, func(t *testing.T) {
			s := fmt.Sprintf(tc.ft, err)
			assert.Regexp(t, tc.expected, s)
		})
	}
}
