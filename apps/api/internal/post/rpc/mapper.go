package rpc

import (
	"github.com/vandad1901/p3s/apps/api/internal/post"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/api/postpb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapToPost(in *postpb.Post) *post.Post {
	return &post.Post{
		ID: in.GetId(),

		Title:  in.GetTitle(),
		Slug:   in.GetSlug(),
		Status: post.PostStatus(in.GetPostStatus()),

		UpdatedAt: in.GetUpdatedAt().AsTime(),
	}
}

func mapToPostBlock(in []*postpb.PostBlock) []*post.PostBlock {
	res := make([]*post.PostBlock, len(in))
	for i, item := range in {
		res[i] = &post.PostBlock{
			ID:     item.GetId(),
			PostID: item.GetPostId(),

			Position:  item.GetPosition(),
			BlockType: post.BlockType(item.GetBlockType()),

			MediaContent: item.GetMedia(),
			TextContent:  item.GetText(),

			Metadata: item.GetMetadata(),
		}
	}

	return res
}

func mapToPostBlockMutateRequest(inserted []*postpb.PostBlock, updated []*postpb.PostBlock, deleted []int64,
) *post.PostBlockMutateRequest {
	return &post.PostBlockMutateRequest{
		Added:   mapToPostBlock(inserted),
		Edited:  mapToPostBlock(updated),
		Removed: deleted,
	}
}

func mapToPostPB(in *post.Post) *postpb.Post {
	return &postpb.Post{
		Id:         in.ID,
		Title:      in.Title,
		Slug:       in.Slug,
		PostStatus: postpb.PostStatus(in.Status),

		CreatedBy: in.CreatedBy,
		UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}

func mapToPostBlockPB(in []*post.PostBlock) []*postpb.PostBlock {
	res := make([]*postpb.PostBlock, len(in))
	for i, item := range in {
		res[i] = &postpb.PostBlock{
			Id:        item.ID,
			Position:  item.Position,
			BlockType: postpb.BlockType(item.BlockType),

			Media:    item.MediaContent,
			Text:     item.TextContent,
			Metadata: item.Metadata,
		}
	}

	return res
}
