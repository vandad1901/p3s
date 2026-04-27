package post

import "purpl3shadow/gen/postpb"

func mapToPost(in *postpb.Post) *Post {
	return &Post{
		ID:       in.GetId(),
		Content:  in.GetContent(),
		AuthorID: in.GetAuthorId(),
	}
}

func mapToPostPB(in *Post) *postpb.Post {
	return &postpb.Post{
		Id:       in.ID,
		Content:  in.Content,
		AuthorId: in.AuthorID,
	}
}
