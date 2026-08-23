package identity

import "time"

type User struct {
	ID int64

	Username     string
	Email        string
	Salt         string
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
}
