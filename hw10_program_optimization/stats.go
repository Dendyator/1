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

func GetDomainStatNew(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	decoder := json.NewDecoder(r)

	for {
		var user User
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}

		atIndex := strings.IndexByte(user.Email, '@')
		emailDomain := user.Email[atIndex+1:]
		if strings.HasSuffix(emailDomain, domain) {
			result[strings.ToLower(emailDomain)]++
		}
	}
	return result, nil
}
