package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	dir := "./testdata/env"

	t.Run("not valid dir", func(t *testing.T) {
		dir := "empty_directory"
		env, err := ReadDir(dir)

		assert.Nil(t, env, "Expected env to be nil")
		assert.Error(t, err, "Expected err not to be nil")
	})

	t.Run("file without value", func(t *testing.T) {
		env, err := ReadDir(dir)

		expectedEnv := Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": EnvValue{
				Value:      "",
				NeedRemove: false,
			},
			"FOO": EnvValue{
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": EnvValue{
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		}

		assert.Equal(t, expectedEnv, env, "Expected env to be equal to expectedEnv")
		assert.Nil(t, err, "Expected err to be nil")
	})
}
