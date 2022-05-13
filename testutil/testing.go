package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func NoError[T any](t *testing.T, f func() (T, error)) T {
	t.Helper()
	r, err := f()
	require.NoError(t, err)
	return r
}
