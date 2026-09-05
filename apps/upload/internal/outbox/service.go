package outbox

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

type Service struct {
	logger    *slog.Logger
	db        *gorm.DB
	publisher *rabbitmq.Publisher
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// MUST BE PASSED CONTEXT TIED TO SYSTEM SIGNALS
func (s *Service) StartWorker(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			err := s.QueueOutbox(ctx)
			if err != nil {
				s.logger.Error("outbox poll failed", "error", err)
			}
		}
	}
}

func (s *Service) QueueOutbox(ctx context.Context) error {
	db := s.db.WithContext(ctx)

	outgoing, claimID, err := getOutgoing(db)
	if err != nil {
		return err
	}

	if len(outgoing) == 0 {
		return nil
	}

	err = send(ctx, db, s.logger,
		s.publisher,
		outgoing, claimID)
	if err != nil {
		return err
	}

	return nil
}

func getOutgoing(db *gorm.DB) ([]Message, string, error) {
	var (
		outgoing []Message
		claimID  string
		err      error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		outgoing, err = dbGetOutgoing(tx)
		if err != nil {
			return err
		}

		claimID, err = dbMarkForTransit(tx, outgoing)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, "", txErr
	}

	return outgoing, claimID, nil
}

func send(ctx context.Context, db *gorm.DB, logger *slog.Logger,
	publisher *rabbitmq.Publisher,
	outgoing []Message, claimID string) error {
	var (
		failed = make([]int64, 0)
	)

	for _, msg := range outgoing {
		err := publisher.PublishWithContext(ctx, msg.MessageBody,
			[]string{msg.RoutingKey},
			rabbitmq.WithPublishOptionsExchange(msg.ExchangeKey))
		if err != nil {
			logger.Error("failed to enqueue message",
				"outbox_id", msg.ID,
				"error", err.Error())
			failed = append(failed, msg.ID)

			if errors.Is(err, rabbitmq.ErrPublishFlowPaused) || errors.Is(err, rabbitmq.ErrPublishBlocked) {
				return fmt.Errorf("queueing message: %w", err)
			}

			continue
		}

		err = db.Where("id = ?", msg.ID).
			Where("in_transit = true").
			Where("claim_id = ?", claimID).
			Delete(&Message{}).Error
		if err != nil {
			return err
		}
	}

	err := db.Model(&Message{}).
		Where("id IN ?", failed).
		Where("in_transit = true").
		Where("claim_id = ?", claimID).
		Updates(map[string]any{
			"in_transit": false,
			"attempts":   gorm.Expr("attempts + 1"),
		}).Error
	if err != nil {
		return err
	}

	return nil
}
