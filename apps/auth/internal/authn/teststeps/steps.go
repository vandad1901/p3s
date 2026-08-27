package teststeps

import (
	"context"
	"strconv"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/cucumber/godog"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vandad1901/p3s/apps/auth/internal/scenario"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/godogutil"
)

func RegisterStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			rowUser := resolveUser(s.SharedData, row)
			password := godogutil.ResolveString(s.SharedData, row, "Password")

			res, err := s.A.AuthnService.Register(ctx, rowUser, password)
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			row["SessionID"] = strconv.FormatInt(res.SessionID, 10)
			row["JWT"] = res.JWT
			row["RefreshToken"] = res.RefreshToken
		}
	}
}

func AssertJWTStep(s *scenario.Scenario) func(context.Context, *godog.Table) {
	return func(ctx context.Context, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			k, err := keyfunc.New(keyfunc.Options{Ctx: ctx, Storage: s.A.KeySet})
			s.Require.NoError(err)

			parsedJWT, err := jwt.ParseWithClaims(godogutil.ResolveString(s.SharedData, row, "JWT"),
				&token.Claims{}, k.Keyfunc, jwt.WithValidMethods([]string{"ES256"}))
			s.Require.NoError(err)
			s.Require.True(parsedJWT.Valid)

			internalClaims, ok := parsedJWT.Claims.(*token.Claims)
			s.Require.True(ok)

			s.Require.NotEmpty(internalClaims.Issuer)
			s.Require.Equal("auth-service", internalClaims.Issuer)
			s.Require.NotEmpty(internalClaims.Subject)
			s.Require.False(internalClaims.IssuedAt.IsZero())
			s.Require.False(internalClaims.ExpiresAt.IsZero())
			s.Require.NotEmpty(internalClaims.ID)

			row["UserID"] = internalClaims.Subject
		}
	}
}

func LoginStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			username := godogutil.ResolveString(s.SharedData, row, "Username")
			password := godogutil.ResolveString(s.SharedData, row, "Password")

			res, err := s.A.AuthnService.Login(ctx, username, password)
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			row["SessionID"] = strconv.FormatInt(res.SessionID, 10)
			row["JWT"] = res.JWT
			row["RefreshToken"] = res.RefreshToken

			s.Require.NotEmpty(res.JWT)
			s.Require.NotEmpty(res.RefreshToken)
			s.Require.True(res.ExpiresAt.After(time.Now()))
		}
	}
}

func RefreshStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			req := resolveRefreshJWTRequest(s, row)

			res, err := s.A.AuthnService.RefreshJWT(ctx, req)
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			s.Require.NotEmpty(res.GetJwt())
			s.Require.NotEqual(row["JWT"], res.GetJwt())

			row["JWT"] = res.GetJwt()
		}
	}
}
