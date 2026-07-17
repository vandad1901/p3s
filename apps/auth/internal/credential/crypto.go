package credential

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func GenerateSalt() ([]byte, string, error) {
	const (
		saltLength = 16
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, "", err
	}

	return salt, base64.RawStdEncoding.EncodeToString(salt), nil
}

func HashPasswordArgon2(password string, salt []byte) string {
	const (
		argonTime    = 1
		argonMemory  = 64 * 1024
		argonThreads = 4
		argonKeyLen  = 32
	)

	hashedPassword := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return base64.RawStdEncoding.EncodeToString(hashedPassword)
}
