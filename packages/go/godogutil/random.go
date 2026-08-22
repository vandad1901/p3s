//go:build test

package godogutil

import "math/rand"

const (
	lowChars  = "abcdefghijklmnopqrstuvwxyz"
	highChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
)

func generateCharFromSet(length int, charSet string) string {
	result := make([]byte, length)
	for i := range length {
		result[i] = charSet[rand.Intn(len(charSet))] //nolint:gosec
	}

	return string(result)
}

func generateRandomChars(length int) string {
	return generateCharFromSet(length, lowChars+highChars)
}

func generateRandomNumbers(length int) string {
	return generateCharFromSet(length, numbers)
}

func generateRandomCharNumbers(length int) string {
	return generateCharFromSet(length, lowChars+highChars+numbers)
}
