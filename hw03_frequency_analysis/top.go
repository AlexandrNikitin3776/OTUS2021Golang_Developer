package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type WordCount struct {
	text  string
	count int
}

const topWordsCount = 10

func Top10(text string) []string {
	re := regexp.MustCompile(`\p{L}+-\p{L}+|\p{L}+`)
	words := re.FindAllString(text, -1)

	frequencyMap := make(map[string]int)
	for _, word := range words {
		frequencyMap[strings.ToLower(word)]++
	}

	wordSlice := make([]WordCount, 0, len(words))
	for word, count := range frequencyMap {
		wordSlice = append(wordSlice, WordCount{word, count})
	}

	sort.Slice(wordSlice, func(i, j int) bool {
		if wordSlice[i].count == wordSlice[j].count {
			return wordSlice[i].text < wordSlice[j].text
		}
		return wordSlice[i].count > wordSlice[j].count
	})

	result := make([]string, 0, 10)
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
