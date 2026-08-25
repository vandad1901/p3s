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

type AuthnRPCServer struct {
	authnpb.UnsafeAuthnServiceServer

	authnService    *authn.Service
	identityService *identity.Service
	sessionService  *session.Service
}

func Register(s *grpc.Server,
	authnService *authn.Service, identityService *identity.Service, sessionService *session.Service) {
	authnpb.RegisterAuthnServiceServer(s, &AuthnRPCServer{
		authnService:    authnService,
		identityService: identityService,
		sessionService:  sessionService,
	})
}

func (s *AuthnRPCServer) Register(ctx context.Context, req *authnpb.RegisterRequest) (*authnpb.RegisterResponse, error) {
	user := mapToUser(req.GetUser())

	res, err := s.authnService.Register(ctx, user, req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("authn register: %w", err)
	}

	return &authnpb.RegisterResponse{Session: mapToSessionResponsePB(res)}, nil
}

func (s *AuthnRPCServer) Login(ctx context.Context, req *authnpb.LoginRequest) (*authnpb.LoginResponse, error) {
	res, err := s.authnService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("authn login: %w", err)
	}

	return &authnpb.LoginResponse{Session: mapToSessionResponsePB(res)}, nil
}

func (s *AuthnRPCServer) RefreshJWT(ctx context.Context, req *authnpb.RefreshJWTRequest,
) (*authnpb.RefreshJWTResponse, error) {
	res, err := s.authnService.RefreshJWT(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("authn refresh jwt: %w", err)
	}

	return res, nil
}
