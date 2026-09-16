package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vandad1901/p3s/apps/media/internal/media"
	"github.com/wagslane/go-rabbitmq"
)

const consumeTimeout = 3 * time.Second

type MediaWorker struct {
	logger   *slog.Logger
	consumer *rabbitmq.Consumer

	mediaService *media.Service
}

func RunLoop(logger *slog.Logger, consumer *rabbitmq.Consumer, mediaService *media.Service) error {
	s := MediaWorker{
		logger:       logger,
		consumer:     consumer,
		mediaService: mediaService,
	}

	err := s.consumer.Run(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			ctx, cancel := context.WithTimeout(context.Background(), consumeTimeout)
			defer cancel()

			err := s.mediaService.HandleMedia(ctx, string(d.Body))
			if err != nil {
				s.logger.Error("handling media",
					"key", string(d.Body),
					"error", err)

				if errors.Is(err, media.ErrAlreadyProcessing) || errors.Is(err, media.ErrAlreadyProcessed) {
					return rabbitmq.Ack
				}

				if errors.Is(err, media.ErrIngestFailed) {
					return rabbitmq.NackDiscard
				}

				return rabbitmq.NackRequeue
			}

			s.logger.Info("media handled successfully",
				"key", string(d.Body))

			return rabbitmq.Ack
		},
	)
	if err != nil {
		return fmt.Errorf("running media consumer loop: %w", err)
	}

	return nil
}
