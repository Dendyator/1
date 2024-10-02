package hw10programoptimization

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	domainSuffix := append([]byte{'.'}, domainBytes...)

	for {
		var user User
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding error: %w", err)
		}

		atIndex := bytes.LastIndexByte([]byte(user.Email), '@')
		if atIndex < 0 || atIndex+1 >= len(user.Email) {
			continue
		}

		emailDomain := []byte(user.Email[atIndex+1:])

		if len(emailDomain) > domainLen && bytes.HasSuffix(emailDomain, domainSuffix) {
			subDomain := emailDomain[:len(emailDomain)-domainLen-1]
			if len(subDomain) > 0 {
				result[string(bytes.ToLower(subDomain))+"."+string(bytes.ToLower(domainBytes))]++
			}
		}
	}
	return result, nil
}
