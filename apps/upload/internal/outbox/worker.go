package outbox

import (
	"context"
	"time"
)

const (
	workerPollInterval = 2 * time.Second
	workerTaskTimeout  = 8 * time.Second
)

func (s *Service) StartWorker() error {
	s.workerStop = make(chan struct{})
	s.workerDone = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	s.workerCancel = cancel

	defer close(s.workerDone)
	defer cancel()

	ticker := time.NewTicker(workerPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.workerStop:
			return nil
		case <-ticker.C:
			taskCtx, taskCancel := context.WithTimeout(ctx, workerTaskTimeout)
			successCount, failedCount, err := s.ProcessQueue(taskCtx)

			taskCancel()

			if err != nil {
				s.logger.Error("outbox poll failed", "error", err)
			} else {
				s.logger.Debug("outbox poll succeeded", "success_count", successCount, "failed_count", failedCount)
			}
		}
	}
}

func (s *Service) GracefulShutdown() {
	s.workerOnce.Do(func() {
		close(s.workerStop)
	})
	<-s.workerDone
}

func (s *Service) ForceShutdown() {
	s.workerOnce.Do(func() {
		close(s.workerStop)
	})
	s.workerCancel()
	<-s.workerDone
}
