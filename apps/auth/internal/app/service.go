package app

import (
	"github.com/vandad1901/p3s/apps/auth/internal/authn"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/jwks"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
)

func initializeServices(a *App, jwtConfig *config.JWTConfig) {
	a.tokenService = token.NewService(a.signer)
	a.identityService = identity.NewService(a.db)
	a.SessionService = session.NewService(a.db, a.tokenService)
	a.AuthnService = authn.NewAuthNService(a.db, a.identityService, a.SessionService, a.tokenService)

	a.JWKSService = jwks.NewService(a.KeySet, jwtConfig.PrivateKey, jwtConfig.KeyID)
}
