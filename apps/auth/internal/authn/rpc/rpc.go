package rpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vandad1901/p3s/apps/auth/internal/authn"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/authnpb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthnRPCServer struct {
	authnpb.UnsafeAuthnServiceServer

	authnService *authn.Service
}

const (
	refreshCookieName     = "refresh_token"
	refreshTokenAgeMonths = 30
	refreshTokenMaxAge    = (time.Hour * 24 * refreshTokenAgeMonths)
)

func Register(s *grpc.Server,
	authnService *authn.Service) {
	authnpb.RegisterAuthnServiceServer(s, &AuthnRPCServer{
		authnService: authnService,
	})
}

func SetRefreshCookie(ctx context.Context, value string) error {
	cookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     "/auth/v1/authn/refresh",
		MaxAge:   int(refreshTokenMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
	if err != nil {
		return fmt.Errorf("setting refresh cookie header: %w", err)
	}

	return nil
}

func (s *AuthnRPCServer) Register(ctx context.Context, req *authnpb.RegisterRequest,
) (*authnpb.RegisterResponse, error) {
	user := mapToUser(req.GetUser())

	res, err := s.authnService.Register(ctx, user, req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("authn register: %w", err)
	}

	err = SetRefreshCookie(ctx, res.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("setting refresh cookie: %w", err)
	}

	return &authnpb.RegisterResponse{Session: mapToSessionResponsePB(res)}, nil
}

func (s *AuthnRPCServer) Login(ctx context.Context, req *authnpb.LoginRequest,
) (*authnpb.LoginResponse, error) {
	res, err := s.authnService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("authn login: %w", err)
	}

	err = SetRefreshCookie(ctx, res.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("setting refresh cookie: %w", err)
	}

	return &authnpb.LoginResponse{Session: mapToSessionResponsePB(res)}, nil
}

var (
	errNoMetadata = errors.New("no incoming metadata")
	errNoCookie   = errors.New("no cookie header present")
	errNoRefresh  = errors.New("refresh cookie not found")
)

func ReadRefreshCookie(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errNoMetadata
	}

	values := md.Get("cookie")
	if len(values) == 0 {
		return "", errNoCookie
	}

	cookies, err := http.ParseCookie(values[0])
	if err != nil {
		return "", fmt.Errorf("parsing cookie header: %w", err)
	}

	for _, c := range cookies {
		if c.Name == refreshCookieName {
			return c.Value, nil
		}
	}

	return "", errNoRefresh
}

func (s *AuthnRPCServer) RefreshJWT(ctx context.Context, req *authnpb.RefreshJWTRequest,
) (*authnpb.RefreshJWTResponse, error) {
	refreshToken, err := ReadRefreshCookie(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading refresh cookie: %w", err)
	}

	res, err := s.authnService.RefreshJWT(ctx, refreshToken, req.GetSessionId(), req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("authn refresh jwt: %w", err)
	}

	return res, nil
}
