package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const buffSize = 512

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	fromFileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	fromFileSize := fromFileStat.Size()
	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer func() {
		if deferErr := fromFile.Close(); deferErr != nil {
			if err == nil {
				err = deferErr
			}
		}
	}()

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func() {
		if deferErr := toFile.Close(); deferErr != nil {
			if err == nil {
				err = deferErr
			}
		}
	}()

	var bytesToRead int64
	if limit > 0 {
		bytesToRead = Min(fromFileSize, limit)
	} else {
		bytesToRead = fromFileSize
	}

	pb := NewProgressBar(bytesToRead)

	if err := copyReaderToWriter(fromFile, fileWriterWithProgressBar, offset, limit); err != nil {
		return err
	}

	return nil
}

func copyReaderToWriter(input io.ReadSeeker, output io.Writer, offset int64, limit int64) error {
	isLimitSet := limit > 0
	var totalBytesCopied int64
	buffBlock := make([]byte, buffSize)

	if _, err := input.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	for (isLimitSet && totalBytesCopied < limit) || !isLimitSet {
		readBytes, err := input.Read(buffBlock)
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		bytesToWrite := int64(readBytes)
		if bytesRemain := limit - totalBytesCopied; isLimitSet && bytesRemain < bytesToWrite {
			bytesToWrite = bytesRemain
		}

		wroteBytes, err := output.Write(buffBlock[:bytesToWrite])
		if err != nil {
			return err
		}
		totalBytesCopied += int64(wroteBytes)
	}
	return nil
}
