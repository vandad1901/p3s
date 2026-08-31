package mediaupload

import "time"

type MediaUploadStatus int32

const (
	MediaUploadStatusUnspecified = iota
	MediaUploadStatusUploading
	MediaUploadStatusUploaded
	MediaUploadStatusIngesting
	MediaUploadStatusIngested
)

type MediaUpload struct {
	ID  int64
	Key string

	status MediaUploadStatus

	CreatedAt time.Time
	UpdatedAt time.Time
	
}

func (*MediaUpload) TableName() string {
	return "media_upload"
}
