package utils

import (
	"crypto/rand"
	"errors"
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

// These are the errors that can be returned by the ValidatePassword function
var ErrPasswordTooShort = errors.New("password too short")
var ErrPasswordTooLong = errors.New("password too long")
var ErrPasswordNoLower = errors.New("password has no lowercase letters")
var ErrPasswordNoUpper = errors.New("password has no uppercase letters")
var ErrPasswordNoNumber = errors.New("password has no numbers")

/*
* ValidatePassword checks if a password meets the following requirements:
* - At least 8 characters long
* - At most 64 characters long
* - Contains at least one lowercase letter
* - Contains at least one uppercase letter
* - Contains at least one number
 */
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if len(password) > 64 {
		return ErrPasswordTooLong
	}

	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return ErrPasswordNoLower
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return ErrPasswordNoUpper
	}

	if !strings.ContainsAny(password, "0123456789") {
		return ErrPasswordNoNumber
	}

	return nil
}
