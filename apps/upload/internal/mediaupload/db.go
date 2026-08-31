package mediaupload

import (
	"context"
	"time"

	"github.com/vandad1901/p3s/packages/go/idv"
	"gorm.io/gorm"
)

func dbCreateWithKey(_ context.Context, db *gorm.DB, key string) (*idv.IDV, error) {
	currentTime := time.Now()

	media := MediaUpload{
		Key: key,

		status:    MediaUploadStatusUploading,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	err := db.Create(media).Error
	if err != nil {
		return nil, err
	}

	return &idv.IDV{ID: media.ID, UpdatedAt: currentTime}, nil
}

func dbChangeStatus(_ context.Context, db *gorm.DB, req *idv.IDV, targetStatus MediaUploadStatus) (*idv.IDV, error) {
	currentTime := time.Now()

	err := db.Model(&MediaUpload{}).
		Where("id = ?", req.ID).
		Where("updated_at = ?", req.UpdatedAt).
		Updates(map[string]any{
			"status": targetStatus,

			"updated_at": currentTime,
		}).Error
	if err != nil {
		return nil, err
	}

	return &idv.IDV{ID: req.ID, UpdatedAt: currentTime}, nil
}

const cutoffHours = 24

func dbDeleteUnused(_ context.Context, db *gorm.DB) error {
	cutoffDate := time.Now().Add(-cutoffHours * time.Hour)

	err := db.Table("media_upload AS mu").
		Joins("LEFT JOIN post_block AS pb ON pb.media_id = mu.id").
		Where("pb.id IS NULL").
		Where("mu.updated_at <= ?", cutoffDate).
		Delete(&MediaUpload{}).Error

	return err
}
