package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func UnpackString(input string) (string, error) {
	var builder strings.Builder
	var lastLetter rune

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

	for i := 0; i < len(part); i++ {
		symbol := part[i]
		if !strings.ContainsRune("0123456789", symbol) {
			builder.WriteRune(symbol)
			lastLetter = symbol
		} else {
			num := int(symbol - '0')
			if num == 0 {
				builder.Reset()
				builder.WriteString(string(part[:i-1]))
			} else {
				builder.WriteString(strings.Repeat(string(lastLetter), num-1))
			}
		}
	}

	return builder.String(), nil
}
