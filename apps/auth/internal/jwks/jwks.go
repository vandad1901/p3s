package jwks

import (
	"context"
	"crypto/ecdsa"
	"log"
	"net/http"

	"github.com/MicahParks/jwkset"
	"github.com/labstack/echo/v4"
)

type Service struct {
	keySet *jwkset.MemoryJWKSet

	privateKey *ecdsa.PrivateKey
	keyID      string
}

func NewService(privateKey *ecdsa.PrivateKey, keyID string) *Service {
	s := &Service{
		keySet: jwkset.NewMemoryStorage(),

		privateKey: privateKey,
		keyID:      keyID,
	}

	err := s.refresh()
	if err != nil {
		log.Fatalf("[!] could not create keyset entry: %s", err)
	}

	return s
}

func (s *Service) Handle(c echo.Context) error {
	response, err := s.keySet.JSONPublic(context.Background())
	if err != nil {
		log.Printf("Failed to get JWK Set JSON: %s", err)
		c.NoContent(http.StatusInternalServerError)

		return nil
	}

	c.JSON(http.StatusOK, response)

	return nil
}

func (s *Service) refresh() error {
	// Maybe implement re-reading jwt_private_key.pem in refresh to handle key rotation
	jwkKey, err := jwkset.NewJWKFromKey(s.privateKey, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{KID: s.keyID},
	})
	if err != nil {
		return err
	}

	s.keySet.KeyWrite(context.Background(), jwkKey)

	return nil
}
