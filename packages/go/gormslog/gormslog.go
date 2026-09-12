package gormslog

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

func New(logger *slog.Logger) *GormSlog {
	return &GormSlog{logger: logger.With("source", "gorm")}
}

type GormSlog struct {
	logger *slog.Logger
}

//nolint:ireturn // required by gorm logger.Interface
func (l GormSlog) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l GormSlog) Info(ctx context.Context, msg string, data ...any) {
	l.logger.InfoContext(ctx, fmt.Sprintf(msg, data...))
}

func (l GormSlog) Warn(ctx context.Context, msg string, data ...any) {
	l.logger.WarnContext(ctx, fmt.Sprintf(msg, data...))
}

func (l GormSlog) Error(ctx context.Context, msg string, data ...any) {
	l.logger.ErrorContext(ctx, fmt.Sprintf(msg, data...))
}

func (l GormSlog) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	sql, rows := fc()

	l.logger.DebugContext(ctx, "gorm query",
		"elapsed", time.Since(begin),
		"rows", rows,
		"sql", sql,
		"error", err,
	)
}
