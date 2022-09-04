package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
)

const buffSize = 512

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	copyFile(from, to, offset, limit)
}

func copyFile(from, to string, offset, limit int64) {
	fromFile, err := os.Open(from)
	if err != nil {
		log.Fatalf("Error while opening %v: %v", from, err)
	}

	defer func() {
		if err := fromFile.Close(); err != nil {
			log.Fatalf("Error while closing file %v: %v", from, err)
		}
	}()

	toFile, err := os.Create(to)
	if err != nil {
		log.Fatalf("Error while opening %v: %v", to, err)
	}

	defer func() {
		if err := toFile.Close(); err != nil {
			log.Fatalf("Error while closing file %v: %v", to, err)
		}
	}()

	if err := copyReaderToWriter(fromFile, toFile, offset, limit); err != nil {
		log.Fatalf("Error while copying from %v to %v: %v", from, to, err)
	}

}

func copyReaderToWriter(input io.ReadSeeker, output io.Writer, offset int64, limit int64) error {
	isLimitSet := limit > 0
	var totalBytesCopied int64
	buffBlock := make([]byte, buffSize)

	if seekedBytes, err := input.Seek(offset, io.SeekStart); err != nil {
		return err
	} else {
		log.Printf("Seeked bytes %v", seekedBytes)
	}

	for (isLimitSet && totalBytesCopied < limit) || !isLimitSet {
		log.Printf("New reading cycle: offset %v, limit %v\n", offset, limit)
		readBytes, err := input.Read(buffBlock)
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		log.Printf("Read bytes %v\n", readBytes)

		bytesToWrite := int64(readBytes)
		if bytesRemain := limit - totalBytesCopied; isLimitSet && bytesRemain < bytesToWrite {
			bytesToWrite = bytesRemain
		}

		wroteBytes, err := output.Write(buffBlock[:bytesToWrite])
		if err != nil {
			return err
		}
		totalBytesCopied += int64(wroteBytes)
		log.Printf("Wrote bytes %v. Total: %v\n", wroteBytes, totalBytesCopied)

	}
	return nil
}
