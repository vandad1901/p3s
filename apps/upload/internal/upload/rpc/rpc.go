package rpc

import (
	"context"
	"fmt"

	"github.com/vandad1901/p3s/apps/upload/internal/upload"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/upload/uploadpb/v1"
	"google.golang.org/grpc"
)

type UploadRPCServer struct {
	uploadpb.UnsafeUploadServiceServer

	uploadService *upload.Service
}

func Register(s *grpc.Server, uploadService *upload.Service) {
	uploadpb.RegisterUploadServiceServer(s, &UploadRPCServer{
		uploadService: uploadService,
	})
}

func (s *UploadRPCServer) GenerateURL(ctx context.Context, req *uploadpb.GenerateURLRequest,
) (*uploadpb.GenerateURLResponse, error) {
	url, fields, err := s.uploadService.GenerateURL(ctx, req.GetMediaKey())
	if err != nil {
		return nil, fmt.Errorf("generating URL: %w", err)
	}

	return &uploadpb.GenerateURLResponse{
		Url:    url,
		Fields: fields,
	}, nil
}

func (s *UploadRPCServer) FinalizeUpload(ctx context.Context, req *uploadpb.FinalizeUploadRequest,
) (*uploadpb.FinalizeUploadResponse, error) {
	err := s.uploadService.FinalizeUpload(ctx, req.GetMediaKey())
	if err != nil {
		return nil, fmt.Errorf("finalizing upload: %w", err)
	}

	return &uploadpb.FinalizeUploadResponse{}, nil
}
