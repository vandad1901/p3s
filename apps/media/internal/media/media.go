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
	Width        int32
	Height       int32

	LeasedAt time.Time
	TryCount int32
}

func (*Media) TableName() string {
	return "media"
}

type derivativeAsset struct {
	file      io.ReadSeeker
	width     int32
	Extension string
}

type Derivative struct {
	MediaID  int64
	MediaKey string

	Ext   string
	Width int32
	URL   string `gorm:"-"`
}

func (*Derivative) TableName() string {
	return "media_derivative"
}

type MediaInfo struct {
	MediaKey    string
	Width       int32
	Height      int32
	Derivatives []Derivative
}
