package post

import (
	"context"
	"purpl3shadow/gen/commonpb"
	"purpl3shadow/gen/postpb"
	"purpl3shadow/utils/dbutil"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Post struct {
	ID       int64
	AuthorID int64
	Content  string

	UpdatedAt time.Time
	CreateAt  time.Time
}

type PostService interface {
	GetPostByID(ctx context.Context, id int64) (*postpb.Post, error)
	CreatePost(ctx context.Context, post *postpb.Post) (*commonpb.IDVersion, error)
	UpdatePost(ctx context.Context, post *postpb.Post) (*commonpb.IDVersion, error)
	DeletePost(ctx context.Context, idv *commonpb.IDVersion) error
}

type PostServiceImpl struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) PostService {
	return &PostServiceImpl{db: db}
}

func (s *PostServiceImpl) GetPostByID(ctx context.Context, id int64) (*postpb.Post, error) {
	var post Post
	if err := s.db.First(&post, id).Error; err != nil {
		return nil, err
	}

	return mapToPostPB(&post), nil
}

func (s *PostServiceImpl) CreatePost(ctx context.Context, post *postpb.Post) (*commonpb.IDVersion, error) {
	idVersion, err := dbCreate(s.db, mapToPost(post))
	if err != nil {
		return nil, err
	}

	return idVersion, nil
}

func (s *PostServiceImpl) UpdatePost(ctx context.Context, in *postpb.Post) (*commonpb.IDVersion, error) {
	var result *commonpb.IDVersion

	err := dbutil.SerializableTx(s.db, func(tx *gorm.DB) error {
		post := mapToPost(in)

		err := tx.Save(&post).Error
		if err != nil {
			return err
		}

		result = &commonpb.IDVersion{
			Id:        post.ID,
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PostServiceImpl) DeletePost(ctx context.Context, idv *commonpb.IDVersion) error {
	return s.db.Delete(&Post{}, idv.GetId()).Error
}
