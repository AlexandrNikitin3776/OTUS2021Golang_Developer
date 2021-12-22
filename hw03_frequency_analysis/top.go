package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const topWordsCount = 10

var regexpWordPattern = regexp.MustCompile(`[\p{L}\d]+[\p{L}\d-]+[\p{L}\d]+|[\p{L}\d]+`)

func getWordsFromText(text string) []string {
	return regexpWordPattern.FindAllString(text, -1)
}

func makeFrequencyMap(words []string) map[string]int {
	if len(words) == 0 {
		return nil
	}
	frequencyMap := make(map[string]int)
	for _, word := range words {
		frequencyMap[strings.ToLower(word)]++
	}
	return frequencyMap
}

func getKeysFromMap(inputMap map[string]int) []string {
	if len(inputMap) == 0 {
		return nil
	}
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}
	return keys
}

func minOfTwo(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func sortSliceByMapValues(sl []string, mp map[string]int) {
	sort.Slice(sl, func(i, j int) bool {
		iSliceValue := mp[sl[i]]
		jSliceValue := mp[sl[j]]
		if iSliceValue == jSliceValue {
			return sl[i] < sl[j]
		}
		return iSliceValue > jSliceValue
	})
}

func Top10(text string) []string {
	words := getWordsFromText(text)
	frequencyMap := makeFrequencyMap(words)
	uniqueWords := getKeysFromMap(frequencyMap)
	sortSliceByMapValues(uniqueWords, frequencyMap)
	return uniqueWords[:minOfTwo(len(uniqueWords), topWordsCount)]
}
