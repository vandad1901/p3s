package auth

import (
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/userpb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapToSessionResponsePB(in *session.SessionResponse) *userpb.AuthSessionResponse {
	return &userpb.AuthSessionResponse{
		SessionId: in.SessionID,

		Jwt:          in.JWT,
		RefreshToken: in.RefreshToken,

		AccessExpiresAt: timestamppb.New(in.ExpiresAt),
	}
}
