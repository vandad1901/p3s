package rpc

import (
	"context"

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

func RegisterUserRPCServer(s *grpc.Server,
	identityService *identity.Service, authService *auth.Service, sessionService *session.Service) {
	userpb.RegisterUserServiceServer(s, &UserRPCServer{
		identityService: identityService,
		authService:     authService,
		sessionService:  sessionService,
	})
}

func (s *UserRPCServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	return s.authService.Register(ctx, req)
}

func (s *UserRPCServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	return s.authService.Login(ctx, req)
}

func (s *UserRPCServer) RefreshJWT(ctx context.Context, req *userpb.RefreshJWTRequest,
) (*userpb.RefreshJWTResponse, error) {
	return s.sessionService.RefreshJWT(ctx, req)
}
