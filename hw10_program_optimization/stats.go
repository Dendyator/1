package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
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

func GetDomainStatOld(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
