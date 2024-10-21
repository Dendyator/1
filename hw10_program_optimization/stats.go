package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go

// easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStatNew(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	var user User
	for scanner.Scan() {
		line := scanner.Bytes()

		if err := user.UnmarshalJSON(line); err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}

		atIndex := strings.IndexByte(user.Email, '@')
		if atIndex < 0 || len(user.Email) <= atIndex+1 {
			continue
		}

		emailDomain := user.Email[atIndex+1:]
		if !strings.HasSuffix(emailDomain, domain) {
			continue
		}

		resultTarget := strings.ToLower(emailDomain)
		result[resultTarget]++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading error: %w", err)
	}

	return result, nil
}
