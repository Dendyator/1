package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var users [100_000]User
	var count int

	decoder := json.NewDecoder(r)
	for {
		if count >= len(users) {
			break
		}
		if err := decoder.Decode(&users[count]); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}
		count++
	}

	domainBytes := []byte(domain)
	domainLen := len(domainBytes)
	domainWithDot := append([]byte{'.'}, domainBytes...)

	for i := 0; i < count; i++ {
		user := users[i]

		atIndex := strings.IndexByte(user.Email, '@')
		if atIndex < 0 {
			continue
		}

		emailDomain := user.Email[atIndex+1:]
		if len(emailDomain) > domainLen && strings.HasSuffix(emailDomain, string(domainWithDot)) {
			subDomain := emailDomain[:len(emailDomain)-domainLen-1]
			if len(subDomain) > 0 {
				result[strings.ToLower(subDomain)+"."+strings.ToLower(domain)]++
			}
		}
	}
	return result, nil
}
