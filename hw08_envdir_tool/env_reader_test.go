package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	dir := "./testdata/env"
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

	t.Run("not valid dir", func(t *testing.T) {
		dir := "empty_directory"
		env, err := ReadDir(dir)

		assert.Nil(t, env, "Expected env to be nil")
		assert.Error(t, err, "Expected err not to be nil")
	})

	t.Run("empty dir", func(t *testing.T) {
		dir := "./empty"
		err := os.Mkdir("./empty", os.ModePerm)
		assert.NoError(t, err)
		defer os.Remove(dir)

		env, err := ReadDir(dir)
		t.Log(err)

		assert.Nil(t, env, "Expected env to be nil")
		assert.Nil(t, err, "Expected err to be nil")
	})

	t.Run("dir instead file", func(t *testing.T) {
	})
	t.Run("test all", func(t *testing.T) {
		err := os.Symlink("./NOTEXISTS", dir+"/ERRFILE")
		assert.NoError(t, err)
		defer os.Remove(dir + "/ERRFILE")

		errDir := dir + "/ERRDIR"
		errD := os.Mkdir(errDir, os.ModePerm)
		assert.NoError(t, errD)
		defer os.Remove(errDir)

		env, err := ReadDir(dir)

		assert.NotNil(t, env)
		assert.Len(t, env, 5)
		assert.Nil(t, err, "Expected err to be nil")
		assert.Equal(t, expectedEnv, env, "Expected env to be equal to expectedEnv")

		os.Setenv("UNSET", "MUST BE REMOVED")
		env.UdpateEnv()

		assert.Equal(t, "bar", os.Getenv("BAR"), "Expected os Env variable BAR=bar")
		_, ok := os.LookupEnv("UNSET")
		assert.False(t, ok, "Expected false")
	})
}

func TestReadFirstLine(t *testing.T) {
	filePath := "testfile.txt"
	fileContent := "Test line\nSecond line"
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Non-existing file", func(t *testing.T) {
		_, err := readFirstLine("non_existing.txt")
		assert.ErrorAs(t, err, &os.ErrNotExist, "Expected os.ErrNotExist")
	})

	t.Run("Empty file", func(t *testing.T) {
		emptyFilePath := "emptyfile.txt"
		emptyFile, err := os.Create(emptyFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(emptyFilePath)
		defer emptyFile.Close()

		result, err := readFirstLine(emptyFilePath)
		assert.NoError(t, err, "Expected no error")
		assert.Zero(t, result, "Expectet zero value of string")
	})

	t.Run("File with multiple lines", func(t *testing.T) {
		expected := "Test line"
		result, err := readFirstLine(filePath)
		assert.NoError(t, err, "Expected no error")
		assert.Equal(t, expected, result, "Expected 'Test line'")
	})
}
