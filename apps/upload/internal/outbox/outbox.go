package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID      int64
	ClaimID uuid.UUID

	ExchangeKey string
	RoutingKey  string
	MessageBody []byte

	QueuedAt    time.Time
	LastTriedAt time.Time
	Attempts    int32
	InTransit   bool
}
