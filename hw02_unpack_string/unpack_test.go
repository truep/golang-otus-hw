package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},

		{input: "а0по5ж", expected: "пооооож"},
		{input: "你3好", expected: "你你你好"},
		{input: "你好4", expected: "你好好好好"},
		{input: "вдф4аыв5жёпы9", expected: "вдффффаывввввжёпыыыыыыыыы"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
		// {input: `qwe\\\\`, expected: `qwe\\`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "১qwe", "ad3১c"}
	invalidNumber := []string{"вl১"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}

	for _, tc := range invalidNumber {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			result, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidNumber), "actual error %q, with result %s", err, result)
		})
	}
}
