package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	text  string
	count int
}

const topWordsCount = 10

func Top10(text string) []string {
	frequencyMap := make(map[string]int)

	words := strings.Fields(text)

	for _, word := range words {
		frequencyMap[word] += 1
	}

	var wordSlice []WordCount
	for word, count := range frequencyMap {
		wordSlice = append(wordSlice, WordCount{word, count})
	}
	sort.Slice(wordSlice, func(i, j int) bool {
		if wordSlice[i].count == wordSlice[j].count {
			return wordSlice[i].text < wordSlice[j].text
		}
		return wordSlice[i].count > wordSlice[j].count
	})

	var result []string
	var upperIndex int

	if len(wordSlice) >= topWordsCount {
		upperIndex = topWordsCount
	} else {
		upperIndex = len(wordSlice)
	}

	for _, word := range wordSlice[:upperIndex] {
		result = append(result, word.text)
	}
	return result
}
