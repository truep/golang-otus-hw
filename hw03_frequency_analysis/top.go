package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	name  string
	count int
}

func Top10(text string) []string {
	if len(text) == 0 {
		return nil
	}

	textFields := strings.Fields(text)
	wordsCounted := make(map[string]int, len(textFields))
	for _, name := range textFields {
		wordsCounted[name]++
	}

	words := make([]Word, 0, len(wordsCounted))
	for name, count := range wordsCounted {
		words = append(words, Word{name, count})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count > words[j].count ||
			(words[i].count == words[j].count &&
				strings.Compare(words[i].name, words[j].name) == -1)
	})

	result := make([]string, 0)
	for _, item := range words {
		result = append(result, item.name)
	}
	if len(result) < 10 {
		return result
	}
	return result[:10]
}
