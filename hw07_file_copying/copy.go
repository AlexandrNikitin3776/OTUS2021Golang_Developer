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
		return fmt.Errorf("error while getting file stat %v: %v", fromPath, err)
	}

	fromFileSize := fromFileStat.Size()
	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error while opening file %v: %v %T", fromPath, err, err)
	}

	defer func() {
		if deferErr := fromFile.Close(); deferErr != nil {
			if err == nil {
				err = fmt.Errorf("error while closing file %v: %v", fromPath, deferErr)
			}
		}
	}()

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error while opening %v: %v", toPath, err)
	}

	defer func() {
		if deferErr := toFile.Close(); deferErr != nil {
			if err == nil {
				err = fmt.Errorf("error while closing file %v: %v", toPath, deferErr)
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

	fileWriterWithProgress := io.MultiWriter(toFile, pb)

	if err := copyReaderToWriter(fromFile, fileWriterWithProgress, offset, limit); err != nil {
		return fmt.Errorf("error while copying from %v to %v: %v", fromPath, toPath, err)
	}

	return nil
}

func copyReaderToWriter(input io.ReadSeeker, output io.Writer, offset int64, limit int64) error {
	isLimitSet := limit > 0
	var totalBytesCopied int64
	buffBlock := make([]byte, buffSize)

	if _, err := input.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("error while seeking file %v", err)
	} else {
	}

	for (isLimitSet && totalBytesCopied < limit) || !isLimitSet {
		readBytes, err := input.Read(buffBlock)
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("error while reading bytes %v", err)
		}

		bytesToWrite := int64(readBytes)
		if bytesRemain := limit - totalBytesCopied; isLimitSet && bytesRemain < bytesToWrite {
			bytesToWrite = bytesRemain
		}

		wroteBytes, err := output.Write(buffBlock[:bytesToWrite])
		if err != nil {
			return fmt.Errorf("error while writing bytes %v", err)
		}
		totalBytesCopied += int64(wroteBytes)
	}
	return nil
}
