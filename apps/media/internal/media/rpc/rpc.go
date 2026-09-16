package rpc

import (
	"context"
	"fmt"

	"github.com/vandad1901/p3s/apps/media/internal/media"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/media/mediapb/v1"
	"google.golang.org/grpc"
)

type MediaRPCServer struct {
	mediapb.UnsafeMediaServiceServer

	mediaService *media.Service
}

func Register(s *grpc.Server, mediaService *media.Service) {
	mediapb.RegisterMediaServiceServer(s, &MediaRPCServer{
		mediaService: mediaService,
	})
}

func (s *MediaRPCServer) MediaIngested(ctx context.Context, req *mediapb.MediaIngestedRequest,
) (*mediapb.MediaIngestedResponse, error) {
	res, err := s.mediaService.MediaIngested(ctx, req.GetMediaKeys())
	if err != nil {
		return nil, fmt.Errorf("checking media are ready: %w", err)
	}

	return &mediapb.MediaIngestedResponse{Unfinished: res}, nil
}
