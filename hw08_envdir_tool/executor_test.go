package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("Empty input", func(t *testing.T) {
		env := make(Environment)
		cmd := []string{}
		returnCode := RunCmd(cmd, env)
		assert.Equal(t, 1, returnCode, "Expected return code 1")
	})

	t.Run("fail command", func(t *testing.T) {
		env := make(Environment)
		cmd := []string{"not_existed_command"}
		returnCode := RunCmd(cmd, env)
		assert.NotEqual(t, 0, returnCode, "Expected non-zero return code")
	})

	t.Run("fail args", func(t *testing.T) {
		env := make(Environment)
		cmd := []string{"ls", "not_exists_args"}
		returnCode := RunCmd(cmd, env)
		assert.NotEqual(t, 0, returnCode, "Expected non-zero return code")
	})

	t.Run("success command with args", func(t *testing.T) {
		env := make(Environment)
		cmd := []string{"echo", "Hello, World!"}
		returnCode := RunCmd(cmd, env)
		assert.Equal(t, 0, returnCode, "Expected return code 0")
	})

	t.Run("success command without args", func(t *testing.T) {
		env := make(Environment)
		cmd := []string{"ls"}
		returnCode := RunCmd(cmd, env)
		assert.Equal(t, 0, returnCode, "Expected return code 0")
	})
}
