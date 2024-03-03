package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	a := strings.Builder{}
	part := []rune(input)
	for k, v := range part {
		if unicode.IsDigit(v) && k == 0 || unicode.IsDigit(v) && unicode.IsDigit(part[k-1]) {
			return "", ErrInvalidString
		}
		if k < len(part)-1 && unicode.IsLetter(v) &&
			!unicode.IsDigit(part[k+1]) || (k == len(part)-1 && unicode.IsLetter(v)) {
			a.WriteRune(v)
		}
		if k < len(part)-1 && unicode.IsSpace(v) && !unicode.IsDigit(part[k+1]) {
			a.WriteRune(v)
		}
		if (unicode.IsDigit(v) && unicode.IsLetter(part[k-1])) || (unicode.IsDigit(v) && unicode.IsSpace(part[k-1])) {
			b := strings.Repeat(string(part[k-1]), int(v)-48)
			for _, e := range b {
				a.WriteRune(e)
			}
		}
	}
	b := a.String()
	a.Reset()
	return b, nil
}
