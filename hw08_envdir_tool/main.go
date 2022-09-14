package main

import (
	"log"
)

func main() {
	_, err := ParseCLIArgs()
	if err != nil {
		log.Fatalf("Can't parse CLI arguments: %v", err)
	}
}
