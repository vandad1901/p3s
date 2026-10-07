package rpc

import (
	"github.com/vandad1901/p3s/apps/media/internal/media"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/media/mediapb/v1"
)

func mapToMediaInfoPBs(mediaInfos []media.MediaInfo) []*mediapb.MediaInfo {
	res := make([]*mediapb.MediaInfo, len(mediaInfos))
	for i, mediaInfo := range mediaInfos {
		res[i] = mapToMediaInfoPB(mediaInfo)
	}

	return res
}

func mapToMediaInfoPB(mediaInfo media.MediaInfo) *mediapb.MediaInfo {
	derivatives := make([]*mediapb.Derivative, len(mediaInfo.Derivatives))
	for i, derivative := range mediaInfo.Derivatives {
		derivatives[i] = &mediapb.Derivative{
			Width: derivative.Width,
			Url:   derivative.URL,
		}
	}

	return &mediapb.MediaInfo{
		MediaKey:    mediaInfo.MediaKey,
		Width:       mediaInfo.Width,
		Height:      mediaInfo.Height,
		Derivatives: derivatives,
	}
}
