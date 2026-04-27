package user

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	CreateUser(username, email, password string) (*User, error)
	UpdateUser(id int, username, email, password string) (*User, error)
	DeleteUser(id int) error
}

type UserServiceImpl struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{db: db}
}

func hashPassword(password string) (string, string, error) {
	const (
		saltLength   = 16
		argonTime    = 1
		argonMemory  = 64 * 1024
		argonThreads = 4
		argonKeyLen  = 32
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hashedPassword), nil
}

func (s *UserServiceImpl) GetUserByID(id int) (*User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) CreateUser(username, email, password string) (*User, error) {
	salt, passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Username:     username,
		Email:        email,
		Salt:         salt,
		PasswordHash: passwordHash,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) UpdateUser(id int, username, email, password string) (*User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	salt, passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email
	user.Salt = salt
	user.PasswordHash = passwordHash
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) DeleteUser(id int) error {
	return s.db.Delete(&User{}, id).Error
}
