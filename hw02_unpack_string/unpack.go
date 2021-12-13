package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const escapeSymbol string = "\\"

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	var resultBuilder strings.Builder
	var targetToRepeat string
	var nextSymbolEscaped bool

	for _, symbolRune := range inputString {
		currentSymbol := string(symbolRune)
		switch {
		case nextSymbolEscaped:
			if !(unicode.IsDigit(symbolRune) || currentSymbol == escapeSymbol) {
				return "", ErrInvalidString
			}
			targetToRepeat = currentSymbol
			nextSymbolEscaped = false

		case currentSymbol == escapeSymbol:
			resultBuilder.WriteString(targetToRepeat)
			targetToRepeat = ""
			nextSymbolEscaped = true

		case unicode.IsDigit(symbolRune):
			if targetToRepeat == "" {
				return "", ErrInvalidString
			}
			repeatCount, _ := strconv.Atoi(currentSymbol)
			resultBuilder.WriteString(strings.Repeat(targetToRepeat, repeatCount))
			targetToRepeat = ""

		default:
			resultBuilder.WriteString(targetToRepeat)
			targetToRepeat = currentSymbol
		}
	}
	if nextSymbolEscaped {
		return "", ErrInvalidString
	}
	resultBuilder.WriteString(targetToRepeat)
	return resultBuilder.String(), nil
}
