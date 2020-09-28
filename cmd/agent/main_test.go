package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		key := []byte(`cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=
    `)
		expected := "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE="
		require.Equal(t, expected, formatKey(key))
	})
}
