package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const buffSize = 512

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error while opening %v: %v", fromPath, err)
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

	if err := copyReaderToWriter(fromFile, toFile, offset, limit); err != nil {
		return fmt.Errorf("error while copying from %v to %v: %v", fromPath, toPath, err)
	}

	return nil
}

func copyReaderToWriter(input io.ReadSeeker, output io.Writer, offset int64, limit int64) error {
	isLimitSet := limit > 0
	var totalBytesCopied int64
	buffBlock := make([]byte, buffSize)

	if seekedBytes, err := input.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("error while seeking file %v", err)
	} else {
		log.Printf("Seeked bytes %v", seekedBytes)
	}

	for (isLimitSet && totalBytesCopied < limit) || !isLimitSet {
		log.Printf("New reading cycle: offset %v, limit %v\n", offset, limit)
		readBytes, err := input.Read(buffBlock)
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("error while reading bytes %v", err)
		}

		log.Printf("Read bytes %v\n", readBytes)

		bytesToWrite := int64(readBytes)
		if bytesRemain := limit - totalBytesCopied; isLimitSet && bytesRemain < bytesToWrite {
			bytesToWrite = bytesRemain
		}

		wroteBytes, err := output.Write(buffBlock[:bytesToWrite])
		if err != nil {
			return fmt.Errorf("error while writing bytes %v", err)
		}
		totalBytesCopied += int64(wroteBytes)

		log.Printf("Wrote bytes %v. Total: %v\n", wroteBytes, totalBytesCopied)

	}
	return nil
}
