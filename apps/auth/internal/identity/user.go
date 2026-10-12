package identity

import "time"

type User struct {
	ID int64

	Username     string
	Email        string
	Salt         string
	PasswordHash string
	DisplayName  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserGetResponse struct {
	ID          int64
	DisplayName string
}
