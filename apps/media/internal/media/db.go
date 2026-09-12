package media

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tryCountLimit = 5
	ingestTimeout = 1 * time.Minute
)

func dbTakeLease(_ context.Context, db *gorm.DB, mediaKey string) (time.Time, error) {
	currentTime := time.Now()
	leaseCutoff := currentTime.Add(-ingestTimeout)

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "media_key"}},
		DoUpdates: clause.Assignments(map[string]any{
			"leased_at": currentTime,
			"try_count": gorm.Expr("media.try_count + 1"),
		}),
		Where: clause.Where{
			Exprs: []clause.Expression{
				gorm.Expr("media.leased_at < ?", leaseCutoff),
				gorm.Expr("media.ingest_status = ?", MediaIngestStatusIngesting),
				gorm.Expr("media.try_count < ?", tryCountLimit),
			},
		},
	}).Create(&Media{
		MediaKey:     mediaKey,
		LeasedAt:     currentTime,
		IngestStatus: MediaIngestStatusIngesting,
		TryCount:     1,
	})

	if result.Error != nil {
		return time.Time{}, result.Error
	}

	if result.RowsAffected > 0 {
		return currentTime, nil
	}

	var media Media

	err := db.
		Where("media_key = ?", mediaKey).
		Select("ingest_status", "try_count").
		Take(&media).Error
	if err != nil {
		return time.Time{}, err
	}

	if media.IngestStatus == MediaIngestStatusIngested {
		return time.Time{}, errAlreadyProcessed
	}

	if media.TryCount >= tryCountLimit {
		return time.Time{}, errIngestFailed
	}

	return time.Time{}, errAlreadyProcessing
}

func dbChangeStatus(_ context.Context, db *gorm.DB,
	mediaKey string, leasedAt time.Time,
	targetStatus MediaIngestStatus) error {
	err := db.Model(&Media{}).
		Where("media_key = ?", mediaKey).
		Where("leased_at = ?", leasedAt).
		Updates(map[string]any{
			"ingest_status": targetStatus,
		}).Error
	if err != nil {
		return err
	}

	return nil
}
