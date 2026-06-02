package authutil

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, string, error) {
	const (
		saltLength   = 16
		argonTime    = 1
		argonMemory  = 64 * 1024
		argonThreads = 4
		argonKeyLen  = 32
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hashedPassword), nil
}
