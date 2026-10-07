package media

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func dbCountReady(db *gorm.DB, keys []string) (int64, error) {
	var count int64

	err := db.Model(&Media{}).
		Where("ingest_status = ?", MediaIngestStatusIngested).
		Where("media_key IN ?", keys).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

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
		return time.Time{}, ErrAlreadyProcessed
	}

	if media.TryCount >= tryCountLimit {
		return time.Time{}, ErrIngestFailed
	}

	return time.Time{}, ErrAlreadyProcessing
}

func dbFinalizeIngest(_ context.Context, db *gorm.DB,
	media *Media) (int64, error) {
	var md Media

	err := db.Model(&Media{}).
		Where("media_key = ?", media.MediaKey).
		Where("leased_at = ?", media.LeasedAt).
		Updates(map[string]any{
			"ingest_status": MediaIngestStatusIngested,
			"width":         media.Width,
			"height":        media.Height,
		}).Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
	).Scan(&md).Error
	if err != nil {
		return 0, err
	}

	return md.ID, nil
}

func dbCreateDerivatives(_ context.Context, db *gorm.DB,
	mediaID int64,
	mediaKey string, derivatives []derivativeAsset) error {
	items := make([]Derivative, len(derivatives))
	for i, derivative := range derivatives {
		items[i] = Derivative{
			MediaID:  mediaID,
			MediaKey: mediaKey,

			Ext:   derivative.Extension,
			Width: derivative.width,
		}
	}

	result := db.Create(&items)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
