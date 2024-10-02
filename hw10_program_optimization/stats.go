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
	decoder := json.NewDecoder(r)

	domainBytes := []byte(domain)
	domainLen := len(domainBytes)
	domainWithDot := append([]byte{'.'}, domainBytes...)

	for {
		var user User
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}

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
