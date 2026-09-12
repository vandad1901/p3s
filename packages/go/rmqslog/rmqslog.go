package rmqslog

import (
	"fmt"
	"log/slog"
)

func New(logger *slog.Logger) *RMQSlog {
	return &RMQSlog{logger: logger.With("source", "rabbitMQ")}
}

type RMQSlog struct {
	logger *slog.Logger
}

func (l RMQSlog) Fatalf(str string, v ...any) {
	l.logger.Error(fmt.Sprintf(str, v...))
}

func (l RMQSlog) Errorf(msg string, data ...any) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

func (l RMQSlog) Warnf(msg string, data ...any) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

func (l RMQSlog) Infof(msg string, data ...any) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

func (l RMQSlog) Debugf(msg string, data ...any) {
	l.logger.Debug(fmt.Sprintf(msg, data...))
}
