package generator

import (
	"math/rand"
	"time"
)

func Password(length int) string {
	if length == 0 {
		length = 10
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	password := make([]byte, length)

	for i := range password {
		password[i] = charset[r.Intn(len(charset))]
	}
	return string(password)
}

func Number(length int) string {
	if length == 0 {
		length = 10
	}
	const charset = "0123456789"
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	password := make([]byte, length)

	for i := range password {
		password[i] = charset[r.Intn(len(charset))]
	}
	return string(password)
}
