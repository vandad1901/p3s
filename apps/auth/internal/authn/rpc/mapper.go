package rpc

import (
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/authnpb/v1"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/userpb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapToUser(in *userpb.User) *identity.User {
	return &identity.User{
		ID:       in.GetId(),
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
	}
}

func mapToSessionResponsePB(in *session.SessionResponse) *authnpb.AuthSessionResponse {
	return &authnpb.AuthSessionResponse{
		SessionId: in.SessionID,

		Jwt:          in.JWT,
		RefreshToken: in.RefreshToken,

		AccessExpiresAt: timestamppb.New(in.ExpiresAt),
	}
}
