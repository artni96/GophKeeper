package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateRandomBytes generates a random byte slice with the given length.
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, err
}
