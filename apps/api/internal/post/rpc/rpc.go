package rpc

import (
	"context"
	"fmt"

	"github.com/vandad1901/p3s/apps/api/internal/post"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/api/postpb/v1"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/commonpb/v1"
	"github.com/vandad1901/p3s/packages/go/idv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCServer struct {
	postpb.UnimplementedPostServiceServer

	postService post.Service
}

func Register(s *grpc.Server, postService *post.Service) {
	postpb.RegisterPostServiceServer(s, &RPCServer{})
}

func (s *RPCServer) Create(ctx context.Context, req *postpb.CreateRequest) (*commonpb.IDVersion, error) {
	p := mapToPost(req.GetPost())
	postBlocks := mapToPostBlock(req.GetPostBlocks())

	res, _, err := s.postService.CreatePost(ctx, p, postBlocks)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return idv.MapToPB(res), nil
}

func (s *RPCServer) Get(ctx context.Context, req *postpb.GetRequest) (*postpb.GetResponse, error) {
	post, postBlocks, err := s.postService.GetPost(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &postpb.GetResponse{
		Post:       mapToPostPB(post),
		PostBlocks: mapToPostBlockPB(postBlocks),
	}, nil
}

func (s *RPCServer) Update(ctx context.Context, req *postpb.UpdateRequest) (*commonpb.IDVersion, error) {
	pst := mapToPost(req.GetPost())

	mutateList := mapToPostBlockMutateRequest(
		req.GetPostBlockRequest().GetInserted(),
		req.GetPostBlockRequest().GetUpdated(),
		req.GetPostBlockRequest().GetDeleted())

	res, _, err := s.postService.UpdatePost(ctx, pst, mutateList)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return idv.MapToPB(res), nil
}

func (s *RPCServer) Delete(ctx context.Context, req *postpb.DeleteRequest) (*emptypb.Empty, error) {
	err := s.postService.DeletePost(ctx, &idv.IDV{})
	if err != nil {
		return nil, fmt.Errorf("error deleting user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
