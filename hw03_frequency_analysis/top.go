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
	wordsCounted := make(map[string]int)
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
				words[i].name < words[j].name)
	})

	result := make([]string, 0)
	maxResult := 10
	for index, item := range words {
		if maxResult > index {
			result = append(result, item.name)
		} else {
			break
		}
	}
	return result
}
