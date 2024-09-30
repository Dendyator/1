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

	for {
		var user User
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}

		emailDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		if strings.HasSuffix(emailDomain, "."+domain) {
			// Извлекаем поддомен вместе с доменом верхнего уровня
			parts := strings.Split(emailDomain, ".")
			if len(parts) >= 2 {
				domainName := strings.Join(parts[:len(parts)-1], ".") + "." + domain
				result[domainName]++
			}
		}
	}
	return result, nil
}
