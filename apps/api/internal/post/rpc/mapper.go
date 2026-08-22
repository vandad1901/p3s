package rpc

import (
	"github.com/vandad1901/p3s/apps/api/internal/post"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/postpb/v1"
)

func mapToPost(in *postpb.Post) *post.Post {
	return &post.Post{
		Title:  in.GetTitle(),
		Slug:   in.GetSlug(),
		Status: post.PostStatus(in.GetPostStatus()),
	}
}

func mapToPostBlock(in []*postpb.PostBlock) []*post.PostBlock {
	res := make([]*post.PostBlock, len(in))
	for i, item := range in {
		res[i] = &post.PostBlock{
			Position:  item.GetPosition(),
			BlockType: post.BlockType(item.GetBlockType()),

			MediaContent: item.GetMedia(),
			TextContent:  item.GetText(),
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
	}
}

func mapToPostBlockPB(in []*post.PostBlock) []*postpb.PostBlock {
	res := make([]*postpb.PostBlock, len(in))
	for i, item := range in {
		res[i] = &postpb.PostBlock{
			Id:        item.ID,
			Position:  item.Position,
			BlockType: postpb.BlockType(item.BlockType),

			Media: item.MediaContent,
			Text:  item.TextContent,
		}
	}

	return res
}
