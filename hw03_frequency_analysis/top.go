package hw03frequencyanalysis

import (
	"cmp"
	"slices"
	"strings"
)

type Words struct {
	word  string
	score int
}

func Top10(text string) []string {
	words := strings.Fields(text)
	scoreMap := make(map[string]int)
	for _, word := range words {
		scoreMap[word]++
	}
	if len(text) == 0 {
		return nil
	}

	topWords := make([]Words, 0, len(scoreMap))
	for k, v := range scoreMap {
		topWords = append(topWords, Words{word: k, score: v})
	}

	slices.SortFunc(topWords, func(a, b Words) int {
		if n := cmp.Compare(b.score, a.score); n != 0 {
			return n
		}
		return cmp.Compare(a.word, b.word)
	})

	if len(topWords) > 10 {
		topWords = topWords[:10]
	}

	var result []string
	for i := 0; i < len(topWords); i++ {
		result = append(result, topWords[i].word)
	}
	return result
}
