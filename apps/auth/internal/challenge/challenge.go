package challenge

import "time"

type Challenge struct {
	ID     int64
	UserID int64

	SecretHash    [64]byte
	ChallengeType string
	Metadata      string

	ExpiresAt  time.Time
	ConsumedAt time.Time
	RevokedAt  time.Time
}


