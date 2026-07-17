package session

import (
	"net"
	"time"
)

type SessionStatus int32

const (
	SessionStatusUnspecified SessionStatus = iota
	SessionStatusActive
	SessionStatusRevoked
)

type Session struct {
	ID     int64
	UserID int64

	RefreshTokenHash string
	IPAddress        net.IP `gorm:"type:inet"`
	UserAgent        string
	Status           SessionStatus

	IssuedAt  time.Time
	ExpiresAt time.Time
}

type SessionResponse struct {
	SessionID    int64
	JWT          string
	RefreshToken string
	ExpiresAt    time.Time
}
