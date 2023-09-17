package xerrors

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errmy = errors.New("plouf")

func TestJoin(t *testing.T) {
	newErr := Join(errmy, "i'm wrapped", errors.New("bla"))

	assert.ErrorIs(t, newErr, errmy)

	var joinErr *joinError
	errors.As(newErr, &joinErr)

	assert.ErrorIs(t, joinErr.err, errmy)
	assert.Len(t, joinErr.errs, 2)
}

func TestJoinError_Error(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")
	for name, test := range map[string]struct {
		err  error
		errs []any
		want string
	}{
		"empty": {
			err:  err1,
			errs: []any{},
			want: "err1",
		},
		"err_ele": {
			err:  err1,
			errs: []any{err2, "a_string"},
			want: "err1: err2 + a_string",
		},
		"nil_elem": {
			err:  err1,
			errs: []any{nil, err2},
			want: "err1: err2",
		},
		"sub_join": {
			err:  err1,
			errs: []any{Join(err2, "sub")},
			want: "err1: err2: sub",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := Join(test.err, test.errs...).Error()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestJoinError_LogValueMethod(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")
	for name, test := range map[string]struct {
		err             error
		errs            []any
		wantMessage     string
		wantStack       string
		wantStackLength int
		wantValues      map[string]any
	}{
		"empty": {
			err:         err1,
			errs:        []any{},
			wantMessage: "err1",
			wantStack: `github.com/emilien-puget/xerrors.TestJoinError_LogValueMethod join_test.go:66
testing.tRunner testing.go:1595
runtime.goexit asm_amd64.s:1650`,
			wantStackLength: 3,
			wantValues:      map[string]any{},
		},
		"err_ele": {
			err:         err1,
			errs:        []any{err2, "a_string"},
			wantMessage: "err1: err2 + a_string",
			wantStack: `github.com/emilien-puget/xerrors.TestJoinError_LogValueMethod join_test.go:66
testing.tRunner testing.go:1595
runtime.goexit asm_amd64.s:1650`,
			wantStackLength: 3,
			wantValues:      map[string]any{},
		},
		"nil_elem": {
			err:         err1,
			errs:        []any{nil, err2},
			wantMessage: "err1: err2",
			wantStack: `github.com/emilien-puget/xerrors.TestJoinError_LogValueMethod join_test.go:66
testing.tRunner testing.go:1595
runtime.goexit asm_amd64.s:1650`,
			wantStackLength: 3,
			wantValues:      map[string]any{},
		},
		"sub_join": {
			err:         err1,
			errs:        []any{Join(err2, "sub")},
			wantMessage: "err1: err2: sub",
			wantStack: `github.com/emilien-puget/xerrors.TestJoinError_LogValueMethod join_test.go:66
testing.tRunner testing.go:1595
runtime.goexit asm_amd64.s:1650`,
			wantStackLength: 3,
			wantValues:      map[string]any{},
		},
		"values": {
			err:         err1,
			errs:        []any{nil, err2, WithValue("one", "two")},
			wantMessage: "err1: err2",
			wantStack: `github.com/emilien-puget/xerrors.TestJoinError_LogValueMethod join_test.go:66
testing.tRunner testing.go:1595
runtime.goexit asm_amd64.s:1650`,
			wantStackLength: 3,
			wantValues:      map[string]any{"one": "two"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			newErr := Join(test.err, test.errs...)
			var buf bytes.Buffer
			h := slog.NewJSONHandler(&buf, nil)
			logger := slog.New(h)
			logger.Info("test", slog.Any("error", newErr))
			var ms []map[string]any
			for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
				if len(line) == 0 {
					continue
				}
				var m map[string]any
				if err := json.Unmarshal(line, &m); err != nil {
					panic(err) // In a real test, use t.Fatal.
				}
				ms = append(ms, m)
			}
			require.Len(t, ms, 1)
			require.Contains(t, ms[0], "error")
			require.IsType(t, map[string]any{}, ms[0]["error"])
			a := ms[0]["error"].(map[string]any)
			assert.Len(t, a, 3)
			assert.NotEmpty(t, a["stacktrace"])
			assert.Equal(t, test.wantMessage, a["message"])
			assert.Equal(t, test.wantValues, a["values"])
		})
	}
}
