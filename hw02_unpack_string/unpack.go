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
	escapeSymbol     = `\`
)

func Unpack(packed string) (string, error) {
	var result strings.Builder
	var skipped bool

	if len(packed) == 0 {
		return "", nil
	}
	if unicode.IsNumber(rune(packed[0])) {
		return "", ErrInvalidString
	}

	end := len(packed) - 1
	for index, width := 0, 0; index < len(packed); index += width {
		curRune, w := utf8.DecodeRuneInString(packed[index:])
		width = w
		if index == end {
			if !unicode.IsNumber(curRune) {
				result.WriteRune(curRune)
			}
			break
		}

		nextRune, _ := utf8.DecodeRuneInString(packed[w+index:])
		switch {
		case (unicode.IsLetter(curRune) || unicode.IsControl(curRune)) && unicode.IsNumber(nextRune):
			repeat, err := getRepeat(nextRune)
			if err != nil {
				return "", ErrInvalidNumber
			}
			result.WriteString(strings.Repeat(string(curRune), repeat))
		case (unicode.IsLetter(curRune) || unicode.IsControl(curRune)) && !unicode.IsNumber(nextRune):
			result.WriteRune(curRune)
		case unicode.IsNumber(curRune) && unicode.IsNumber(nextRune) && !skipped:
			return "", ErrInvalidString
		case (checkEsc(curRune) && checkEsc(nextRune) && !skipped):
			result.WriteRune(nextRune)
			skipped = true
		case (checkEsc(curRune) && unicode.IsLetter(nextRune)):
			return "", ErrInvalidString
		case (checkEsc(curRune) && unicode.IsNumber(nextRune) && !skipped):
			result.WriteRune(nextRune)
			skipped = true
		case (checkEsc(curRune) || unicode.IsNumber(curRune)) && unicode.IsNumber(nextRune):
			repeat, err := getRepeat(nextRune)
			if err != nil {
				return "", ErrInvalidNumber
			}
			result.WriteString(strings.Repeat(string(curRune), repeat-1))
		default:
			skipped = false
		}
	}
	return result.String(), nil
}

func getRepeat(next rune) (int, error) {
	repeat, err := strconv.Atoi(string(next))
	if err != nil {
		return 0, ErrInvalidNumber
	}
	return repeat, nil
}

func checkEsc(r rune) bool {
	return string(r) == escapeSymbol
}
