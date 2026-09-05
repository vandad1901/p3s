package outbox

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	batchSize     = 20
	attemptsLimit = 3
)

func dbGetOutgoing(db *gorm.DB) ([]Message, error) {
	var outgoing []Message

	err := db.
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
		// recovered messages skip the attempt limit
		Where(`
			(in_transit = false AND attempts < ?)
			OR 
			(in_transit = true AND last_tried_at <= ?)
			`, attemptsLimit, time.Now().Add(-15*time.Minute)).
		Order("queued_at").
		Limit(batchSize).
		Find(&outgoing).Error
	if err != nil {
		return nil, err
	}

	return outgoing, nil
}

func dbMarkForTransit(db *gorm.DB, messages []Message) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}

	ids := make([]int64, len(messages))

	for i, msg := range messages {
		ids[i] = msg.ID
	}

	claimID := uuid.NewString()

	err := db.Model(&Message{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"claim_id": claimID,

			"in_transit":    true,
			"last_tried_at": time.Now(),
		}).Error
	if err != nil {
		return "", err
	}

	return claimID, nil
}
