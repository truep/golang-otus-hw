package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	srcPath := "/tmp/srcFile.txt"
	dstPath := "/tmp/dstFile.txt"
	srcFile, _ := os.Create(srcPath)
	srcFile.WriteString("Test data")
	srcFile.Close()
	defer os.Remove(srcPath)
	defer os.Remove(dstPath)

	t.Run("Successful copy", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 0, 0)
		assert.NoError(t, err)

		srcStat, _ := os.Stat(srcPath)
		stat, err := os.Stat(dstPath)
		assert.NoError(t, err)
		assert.EqualValues(t, srcStat.Size(), stat.Size())
	})

	t.Run("Equal from and to path", func(t *testing.T) {
		err := Copy(srcPath, srcPath, 0, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrEqualPath, err)
	})

	t.Run("Unsupported src file", func(t *testing.T) {
		err := Copy("/dev/urandom", dstPath, 0, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("Offset > file size", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 100, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("Copy Invalid file path", func(t *testing.T) {
		src := "/invalid/path"
		err := Copy(src, dstPath, 0, 0)
		assert.Error(t, err)
	})

	t.Run("Copy to Invalid file path", func(t *testing.T) {
		dst := "/invalid/path"
		err := Copy(srcPath, dst, 0, 0)
		assert.Error(t, err)
	})

	t.Run("Copy with offset and limit", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 0, 5)
		assert.NoError(t, err)

		stat, err := os.Stat(dstPath)
		assert.NoError(t, err)
		assert.EqualValues(t, 5, stat.Size())
	})

	t.Run("Copy with limit > file size", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 0, 100)
		assert.NoError(t, err)

		stat, err := os.Stat(dstPath)
		assert.NoError(t, err)
		assert.EqualValues(t, 9, stat.Size())
	})
	t.Run("Copy with non-zero offset and limit", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 2, 2)
		assert.NoError(t, err)

		stat, err := os.Stat(dstPath)
		assert.NoError(t, err)
		assert.EqualValues(t, 2, stat.Size())

		dstFile, _ := os.Open(dstPath)
		defer dstFile.Close()

		dstContent := make([]byte, 2)
		dstFile.Read(dstContent)
		assert.Equal(t, []byte("st"), dstContent)
	})

	t.Run("Copy with zero offset and non-zero limit", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 0, 3)
		assert.NoError(t, err)

		stat, err := os.Stat(dstPath)
		assert.NoError(t, err)
		assert.EqualValues(t, 3, stat.Size())

		dstFile, _ := os.Open(dstPath)
		defer dstFile.Close()

		dstContent := make([]byte, 3)
		dstFile.Read(dstContent)
		assert.Equal(t, []byte("Tes"), dstContent)
	})

	t.Run("Copy with non-zero offset and limit > file size", func(t *testing.T) {
		err := Copy(srcPath, dstPath, 5, 100)
		assert.NoError(t, err)

		stat, err := os.Stat(dstPath)
		assert.NoError(t, err)
		assert.EqualValues(t, 4, stat.Size())

		dstFile, _ := os.Open(dstPath)
		defer dstFile.Close()

		dstContent := make([]byte, 4)
		dstFile.Read(dstContent)
		assert.Equal(t, []byte("data"), dstContent)
	})
}
