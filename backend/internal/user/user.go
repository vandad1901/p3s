package user

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Salt         string `json:"salt"`
	PasswordHash string `json:"password_hash"`
}
