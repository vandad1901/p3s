package idv

import (
	"time"

	"github.com/vandad1901/p3s/packages/go/gen/protobuf/commonpb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IDV struct {
	ID        int64
	UpdatedAt time.Time
}

func MapToPB(in *IDV) *commonpb.IDVersion {
	return &commonpb.IDVersion{
		Id:        in.ID,
		UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}
