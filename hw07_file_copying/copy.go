package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEqualPath             = errors.New("from and to path must not be equal")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrEqualPath
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fstat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if !fstat.Mode().IsRegular() || fstat.IsDir() {
		return ErrUnsupportedFile
	}

	if fstat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		dstFile.Close()
		if err != nil {
			os.Remove(toPath)
		}
	}()

	if limit == 0 {
		limit = fstat.Size()
	}

	wSize := getDstSize(fstat.Size(), offset, limit)
	section := io.NewSectionReader(srcFile, offset, limit)

	bar := pb.Full.Start64(wSize)
	defer bar.Finish()

	barReader := bar.NewProxyReader(section)

	_, err = io.Copy(dstFile, barReader)
	if err != nil {
		return err
	}
	return nil
}

func getDstSize(srcSize, offset, limit int64) int64 {
	if limit >= (srcSize-offset) || limit == 0 {
		return srcSize - offset
	}
	return limit
}
