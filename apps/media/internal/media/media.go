package media

import (
	"io"
	"time"
)

type MediaIngestStatus int32

const (
	MediaIngestStatusUnspecified = iota
	MediaIngestStatusIngesting
	MediaIngestStatusIngested
)

type Media struct {
	ID       int64
	MediaKey string

	IngestStatus MediaIngestStatus
	MimeType     string
	Size         int64

	LeasedAt time.Time
	TryCount int32
}

func (*Media) TableName() string {
	return "media"
}

type derivative struct {
	file      io.ReadSeeker
	Extension string
}
