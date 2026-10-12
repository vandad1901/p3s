package rpc

import (
	"context"
	"fmt"

	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/userpb/v1"
	"google.golang.org/grpc"
)

type RPCServer struct {
	userpb.UnsafeUserServiceServer

	identityService *identity.Service
}

func Register(s *grpc.Server, identityService *identity.Service) {
	userpb.RegisterUserServiceServer(s, &RPCServer{
		identityService: identityService,
	})
}

func (s *RPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	res, err := s.identityService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}

	return &userpb.GetUserResponse{
		User: mapToUserPB(res),
	}, nil
}

func (s *RPCServer) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	users, err := s.identityService.GetUsers(ctx, req.GetUserIds())
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	var userResponses []*userpb.UserGetResponse
	for _, user := range users {
		userResponses = append(userResponses, &userpb.UserGetResponse{
			Id:          user.ID,
			DisplayName: user.DisplayName,
		})
	}

	return &userpb.GetUsersResponse{
		Users: userResponses,
	}, nil
}

func (s *RPCServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	err := s.identityService.UpdateUser(ctx, &identity.User{
		ID:          req.GetId(),
		DisplayName: req.GetDisplayName(),
	})
	if err != nil {
		return nil, fmt.Errorf("error updating user info: %w", err)
	}

	return &userpb.UpdateUserResponse{}, nil
}
