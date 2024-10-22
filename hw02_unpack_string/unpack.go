package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func UnpackString(input string) (string, error) {
	var letter []string
	var lastLetter string
	for i := 10; i <= 99; i++ {
		re := regexp.MustCompile(strconv.Itoa(i))
		if re.MatchString(input) {
			return "", ErrInvalidString
		}
	}
	part := []rune(input)
	if len(part) == 0 {
		return "", nil
	}
	if strings.ContainsRune("0123456789", part[0]) {
		return "", ErrInvalidString
	}
	for _, symbol := range part {
		if !strings.ContainsRune("0123456789", symbol) {
			letter = append(letter, string(symbol))
			lastLetter = string(symbol)
		} else {
			letter = remove(letter)
			num := strings.IndexRune("0123456789", symbol)
			letter = append(letter, strings.Repeat(lastLetter, num))
		}
	}
	a := strings.Join(letter, "")
	return a, nil
}

func remove(s []string) []string {
	return s[:len(s)-1]
}
