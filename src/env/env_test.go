package env_test

import (
	"testing"

	"github.com/buildkite/ecrscanresults/src/env"
	"github.com/stretchr/testify/assert"
)

func TestParseWithPrefix(t *testing.T) {
	t.Run("read when no env is set", func(t *testing.T) {
		got := env.ParseWithPrefix("NOT_PRESENT")
		want := []string{}
		assert.Equal(t, want, got)
	})
	t.Run("read multiple env vars", func(t *testing.T) {
		t.Setenv("SOME_VAR_0", "val1")
		t.Setenv("SOME_VAR_1", "val2")
		t.Setenv("SOME_VAR_3", "")
		got := env.ParseWithPrefix("SOME_VAR_")
		want := []string{"val1", "val2"}
		assert.Equal(t, want, got)
	})
}
