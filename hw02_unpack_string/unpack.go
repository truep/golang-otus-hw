package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidNumber = errors.New("invalid number")
	ErrInvalidString = errors.New("invalid string")
)

func Unpack(packed string) (string, error) {
	if len(packed) == 0 {
		return "", nil
	}

	if !firstOk(packed) {
		return "", ErrInvalidString
	}

	var result strings.Builder

	for len(packed) > 0 {
		curRune, size := utf8.DecodeRuneInString(packed)
		nextRune, _ := utf8.DecodeRuneInString(packed[size:])
		if unicode.IsLetter(curRune) || unicode.IsControl(curRune) {
			if unicode.IsNumber(nextRune) {
				repeat, err := strconv.Atoi(string(nextRune))
				if err != nil {
					return "", ErrInvalidNumber
				}
				result.WriteString(strings.Repeat(string(curRune), repeat))
			} else {
				result.WriteRune(curRune)
			}
		}
		if unicode.IsNumber(curRune) && unicode.IsNumber(nextRune) {
			return "", ErrInvalidString
		}
		packed = packed[size:]
	}
	return result.String(), nil
}

func firstOk(s string) bool {
	v, _ := utf8.DecodeRuneInString(s[0:])
	return !unicode.IsNumber(v)
}
