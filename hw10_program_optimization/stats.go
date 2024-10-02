package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
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

	start := time.Now()
	defer func() {
		fmt.Printf("Execution time: %v\n", time.Since(start))
	}()

	for {
		var user User
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}

		atIndex := strings.IndexByte(user.Email, '@')
		if atIndex < 0 || atIndex+1 >= len(user.Email) {
			continue
		}

		emailDomain := strings.ToLower(user.Email[atIndex+1:])

		if len(emailDomain) > len(domain) && strings.HasSuffix(emailDomain, "."+domain) {
			subDomain := emailDomain[:len(emailDomain)-len(domain)-1]
			if subDomain != "" {
				result[subDomain+"."+domain]++
			}
		}
	}
	return result, nil
}
