package rpc

import (
	"context"
	"fmt"

	"github.com/vandad1901/p3s/apps/auth/internal/authn"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/authnpb/v1"
	"google.golang.org/grpc"
)

type AuthRPCServer struct {
	authnpb.UnsafeAuthnServiceServer

	identityService *identity.Service
	authService     *authn.Service
	sessionService  *session.Service
}

func Register(s *grpc.Server,
	identityService *identity.Service, authService *authn.Service, sessionService *session.Service) {
	authnpb.RegisterAuthnServiceServer(s, &AuthRPCServer{
		identityService: identityService,
		authService:     authService,
		sessionService:  sessionService,
	})
}

func (s *AuthRPCServer) Register(ctx context.Context, req *authnpb.RegisterRequest) (*authnpb.RegisterResponse, error) {
	user := mapToUser(req.GetUser())

	res, err := s.authService.Register(ctx, user, req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("error while registering user: %w", err)
	}

	return &authnpb.RegisterResponse{Session: mapToSessionResponsePB(res)}, nil
}

func (s *AuthRPCServer) Login(ctx context.Context, req *authnpb.LoginRequest) (*authnpb.LoginResponse, error) {
	res, err := s.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("error while logging in user: %w", err)
	}

	return &authnpb.LoginResponse{Session: mapToSessionResponsePB(res)}, nil
}

func (s *AuthRPCServer) RefreshJWT(ctx context.Context, req *authnpb.RefreshJWTRequest,
) (*authnpb.RefreshJWTResponse, error) {
	res, err := s.sessionService.RefreshJWT(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while refreshing JWT token user: %w", err)
	}

	return res, nil
}
