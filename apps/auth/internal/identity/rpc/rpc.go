package rpc

import (
	"context"
	"fmt"

	"github.com/vandad1901/p3s/apps/auth/internal/auth"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	userpb "github.com/vandad1901/p3s/packages/go/gen/protobuf/userpb/v1"
	"google.golang.org/grpc"
)

type UserRPCServer struct {
	userpb.UnsafeUserServiceServer

	identityService *identity.Service
	authService     *auth.Service
	sessionService  *session.Service
}

func Register(s *grpc.Server,
	identityService *identity.Service, authService *auth.Service, sessionService *session.Service) {
	userpb.RegisterUserServiceServer(s, &UserRPCServer{
		identityService: identityService,
		authService:     authService,
		sessionService:  sessionService,
	})
}

func (s *UserRPCServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	res, err := s.authService.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while registering user: %w", err)
	}

	return res, nil
}

func (s *UserRPCServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	res, err := s.authService.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while logging in user: %w", err)
	}

	return res, nil
}

func (s *UserRPCServer) RefreshJWT(ctx context.Context, req *userpb.RefreshJWTRequest,
) (*userpb.RefreshJWTResponse, error) {
	res, err := s.sessionService.RefreshJWT(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while refreshing JWT token user: %w", err)
	}

	return res, nil
}
