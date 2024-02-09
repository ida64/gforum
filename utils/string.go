package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

/*
* GenerateRandomString generates a random string of the given length
 */
func GenerateRandomString(length int) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var builder strings.Builder
	builder.Grow(length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphanumeric))))
		builder.WriteByte(alphanumeric[randomIndex.Int64()])
	}

	return builder.String()
}

func StringArrayContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
