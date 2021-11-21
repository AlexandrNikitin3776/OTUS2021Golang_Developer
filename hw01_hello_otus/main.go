package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	inputString := "Hello, OTUS!"
	reversedString := stringutil.Reverse(inputString)
	fmt.Println(reversedString)
}
