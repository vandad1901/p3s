package jwks

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MicahParks/jwkset"
)

type Service struct {
	keySet *jwkset.MemoryJWKSet

	privateKey *ecdsa.PrivateKey
	keyID      string
}

func NewService(keyset *jwkset.MemoryJWKSet, privateKey *ecdsa.PrivateKey, keyID string) *Service {
	s := &Service{
		keySet: keyset,

		privateKey: privateKey,
		keyID:      keyID,
	}

	err := s.refresh()
	if err != nil {
		log.Fatalf("[!] could not create keyset entry: %s", err)
	}

	return s
}

func (s *Service) Handle(ctx context.Context) (json.RawMessage, error) {
	response, err := s.keySet.JSONPublic(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting public keyset: %w", err)
	}

	return response, nil
}

func (s *Service) refresh() error {
	// Maybe implement re-reading jwt_private_key.pem in refresh to handle key rotation
	jwkKey, err := jwkset.NewJWKFromKey(s.privateKey, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{KID: s.keyID},
	})
	if err != nil {
		return fmt.Errorf("creating JWK from private key: %w", err)
	}

	err = s.keySet.KeyWrite(context.Background(), jwkKey)
	if err != nil {
		return fmt.Errorf("writing JWK to keyset: %w", err)
	}

	return nil
}
