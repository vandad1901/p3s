package credential

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func GenerateSalt() ([]byte, string, error) {
	const (
		saltLength = 16
	)

	salt := make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, "", fmt.Errorf("generating salt: %w", err)
	}

	return salt, base64.RawStdEncoding.EncodeToString(salt), nil
}

func HashPasswordArgon2(password string, salt []byte) string {
	const (
		argonTime    = 2
		argonMemory  = 19 * 1024
		argonThreads = 1
		argonKeyLen  = 32
	)

	hashedPassword := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return base64.RawStdEncoding.EncodeToString(hashedPassword)
}
