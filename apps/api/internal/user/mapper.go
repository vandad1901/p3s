package user

import (
	"purpl3shadow/gen/userpb"
	"purpl3shadow/utils/authutil"
)

func mapToUser(in *userpb.User) (*User, error) {
	salt, passwordHash, err := authutil.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           in.Id,
		Username:     in.Username,
		Email:        in.Email,
		Salt:         salt,
		PasswordHash: passwordHash,
	}, nil
}

func mapToUserPB(user *User) *userpb.User {
	return &userpb.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		
	}
}
