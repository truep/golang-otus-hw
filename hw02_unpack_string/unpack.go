package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidNumber = errors.New("invalid number")
	ErrInvalidString = errors.New("invalid string")
)

func Unpack(packed string) (string, error) {
	var result strings.Builder
	if len(packed) == 0 {
		return "", nil
	}
	runePacked := []rune(packed)
	if unicode.IsNumber(runePacked[0]) {
		return "", ErrInvalidString
	}

	for index, curRune := range runePacked {
		if index == len(runePacked)-1 {
			if !unicode.IsNumber(curRune) {
				result.WriteRune(curRune)
			}
			break
		}
		nextRune := runePacked[index+1]
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
	}
	return result.String(), nil
}
