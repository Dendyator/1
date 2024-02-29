package hw03frequencyanalysis

import (
	"sort"
	"strconv"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)
	scoreMap := make(map[string]int)
	for _, word := range words {
		scoreMap[word]++
	}
	if len(text) == 0 {
		return nil
	}
	sortedNumbers := make([][2]string, 0, len(scoreMap))
	for word, number := range scoreMap {
		sortedNumbers = append(sortedNumbers, [2]string{word, strconv.Itoa(number)})
	}
	sort.Slice(sortedNumbers, func(i, j int) bool {
		return sortedNumbers[i][1] > sortedNumbers[j][1] || (sortedNumbers[i][1] == sortedNumbers[j][1]) && sortedNumbers[i][0] < sortedNumbers[j][0]
	})
	sortedNumbers = sortedNumbers[:10]
	result := []string{}
	for i := 0; i < 10; i++ {
		word := sortedNumbers[i][0]
		result = append(result, word)
	}
	return result
}
