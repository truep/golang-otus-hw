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
	wordCounted := []Word{}

	if len(text) == 0 {
		return nil
	}

	textFields := strings.Fields(text)

	for _, name := range textFields {
		found := false
		for i := 0; i < len(wordCounted); i++ {
			if wordCounted[i].name == name {
				wordCounted[i].count++
				found = true
				break
			}
		}
		if !found {
			wordCounted = append(wordCounted, Word{name, 1})
		}
	}
	sort.Slice(wordCounted, func(i, j int) bool {
		return wordCounted[i].count > wordCounted[j].count ||
			(wordCounted[i].count == wordCounted[j].count &&
				strings.Compare(wordCounted[i].name, wordCounted[j].name) == -1)
	})

	result := make([]string, 0)
	for _, item := range wordCounted {
		result = append(result, item.name)
	}
	if len(result) < 10 {
		return result
	}
	return result[:10]
}
