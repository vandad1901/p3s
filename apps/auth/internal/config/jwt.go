package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/vandad1901/p3s/packages/go/envutil"
)

type JWTConfig struct {
	PrivateKey *ecdsa.PrivateKey
	KeyID      string
}

func loadJWTConfig() *JWTConfig {
	ecKey := MustGetECKey()
	KeyID := envutil.MustGetString("KEYSET_KEY_ID")

	return &JWTConfig{
		PrivateKey: ecKey,
		KeyID:      KeyID,
	}
}

func MustGetECKey() *ecdsa.PrivateKey {
	privateKeyPath := envutil.MustGetString("JWT_PRIVATE_KEY_PATH")

	privateKeyPEM, err := os.ReadFile(privateKeyPath) //nolint:gosec
	if err != nil {
		log.Fatalf("Invalid JWT_PRIVATE_KEY_PATH. Unable to read file: %s", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatalf("Invalid JWT_PRIVATE_KEY_PATH. Private key must be in PEM format")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Invalid JWT_PRIVATE_KEY_PATH. Unable to parse private key: %s", err)
	}

	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		log.Fatalf("Invalid JWT_PRIVATE_KEY_PATH. Not an ECDSA private key")
	}

	return ecKey
}
