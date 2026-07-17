package identity

import (
	userpb "github.com/vandad1901/p3s/packages/go/gen/protobuf/userpb/v1"
)

func MapToUser(in *userpb.User) (*User, error) {
	return &User{
		ID:       in.GetId(),
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
	}, nil
}
