package rpc

import (
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/userpb/v1"
)

func mapToUserPB(user *identity.User) *userpb.User {
	return &userpb.User{
		Id:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}
}
