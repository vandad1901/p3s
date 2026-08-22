package post

import (
	"time"

	"github.com/vandad1901/p3s/packages/go/mutatelist"
)

type PostStatus int32

const (
	PostStatusUnspecified PostStatus = iota
	PostStatusDraft
	PostStatusPublished
)

type Post struct {
	ID int64

	Title  string
	Slug   string
	Status PostStatus

	CreatedAt time.Time
	CreatedBy int64
	UpdatedAt time.Time
	UpdatedBy int64
}

func (*Post) TableName() string {
	return "post"
}

type BlockType int32

const (
	BlockTypeUnspecified BlockType = iota
	BlockTypeText
	BlockTypeMedia
)

type PostBlock struct {
	ID       int64
	PostID   int64
	Position int32

	BlockType    BlockType
	MediaContent string
	TextContent  string
	Metadata     string
}

func (*PostBlock) TableName() string {
	return "post_block"
}

type PostBlockMutateRequest = mutatelist.MutateList[*PostBlock, int64]

func (pb *PostBlock) GetID() int64 {
	return pb.ID
}

func (pb *PostBlock) GetPosition() int {
	return int(pb.Position)
}
