package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func remove(s []string) []string {
	return s[:len(s)-1]
}

func Unpack(input string) (string, error) {
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
	if strings.IndexRune("0123456789", part[0]) != -1 { //nolint:gosimple
		return "", ErrInvalidString
	}
	for _, symbol := range part {
		if strings.IndexRune("0123456789", symbol) == -1 { //nolint:gosimple
			letter = append(letter, string(symbol))
			lastLetter = string(symbol)
		} else {
			letter = remove(letter)
			num := strings.IndexRune("0123456789", symbol)
			letter = append(letter, strings.Repeat(lastLetter, num))
		}
	}
	a := strings.Join(letter, "")
	letter = []string{}
	return a, nil
}
