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
	a := make([][2]string, 0, len(scoreMap))
	for word, number := range scoreMap {
		a = append(a, [2]string{word, strconv.Itoa(number)})
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i][1] > a[j][1] || (a[i][1] == a[j][1]) && a[i][0] < a[j][0]
	})
	a = a[:10]
	var result []string
	for i := 0; i < 10; i++ {
		word := a[i][0]
		result = append(result, word)
	}
	return result
}
