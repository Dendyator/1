package hw10programoptimization

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"testing"
)

func generateTestUserData(userID int) io.Reader {
	user := User{
		ID:       userID,
		Name:     "User " + strconv.Itoa(userID),
		Email:    "user" + strconv.Itoa(userID) + "@example.com",
		Username: "username" + strconv.Itoa(userID),
		Phone:    "123-456-7890",
		Password: "password",
		Address:  "Address",
	}

	data, _ := json.Marshal(user)
	return bytes.NewReader(data)
}

func BenchmarkGetDomainStatNew(b *testing.B) {
	const numEntries = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 1; j <= numEntries; j++ {
			data := generateTestUserData(j)
			_, err := GetDomainStatNew(data, "com")
			if err != nil {
				b.Errorf("Error: %v", err)
			}
		}
	}
}
func BenchmarkGetDomainStatOld(b *testing.B) {
	const numEntries = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 1; j <= numEntries; j++ {
			data := generateTestUserData(j)
			_, err := GetDomainStatOld(data, "com")
			if err != nil {
				b.Errorf("Error: %v", err)
			}
		}
	}
}
