package main

import (
	"flag"
	"io"
	"log"
	"os"
)

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
	dd(from, to, limit, offset)
}

func dd(from, to string, limit, offset int64) {
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

	if written, err := io.Copy(toFile, fromFile); err != nil {
		log.Fatalf("Error while copying from %v to %v: %v", from, to, err)
	} else {
		log.Printf("Copied %v bytes", written)
	}

}
