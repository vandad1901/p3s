package post

import (
	"context"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/idv"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreatePost(ctx context.Context, post *Post, postBlocks []*PostBlock,
) (*idv.IDV, []int64, error) {
	db := s.db.WithContext(ctx)

	err := validatePost(ctx, db, post, postBlocks)
	if err != nil {
		return nil, nil, err
	}

	var (
		res      *idv.IDV
		addedIDs []int64
	)

	err = dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbCreatePost(ctx, tx, post)
		if err != nil {
			return err
		}

		addedIDs, err = dbCreatePostBlocks(ctx, tx, res.ID, postBlocks)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return res, addedIDs, nil
}

func (s *Service) GetPost(ctx context.Context, postID int64) (*Post, []*PostBlock, error) {
	db := s.db.WithContext(ctx)

	var (
		post       *Post
		postBlocks []*PostBlock
		err        error
	)

	err = dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		post, postBlocks, err = GetTx(ctx, tx, postID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return post, postBlocks, nil
}

func GetTx(ctx context.Context, tx *gorm.DB, postID int64) (*Post, []*PostBlock, error) {
	post, err := dbGetPost(ctx, tx, postID)
	if err != nil {
		return nil, nil, err
	}

	postBlocks, err := dbGetPostBlocks(ctx, tx, postID)
	if err != nil {
		return nil, nil, err
	}

	return post, postBlocks, nil
}

func (s *Service) UpdatePost(ctx context.Context, post *Post, mutateRequest *PostBlockMutateRequest,
) (*idv.IDV, []int64, error) {
	db := s.db.WithContext(ctx)

	err := validatePostForUpdate(ctx, db, post, mutateRequest)
	if err != nil {
		return nil, nil, err
	}

	var (
		res      *idv.IDV
		addedIDs []int64
	)

	err = dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, addedIDs, err = updatePostTx(ctx, tx, post, mutateRequest)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return res, addedIDs, nil
}

func updatePostTx(ctx context.Context, tx *gorm.DB, post *Post, mutateRequest *PostBlockMutateRequest,
) (*idv.IDV, []int64, error) {
	res, err := dbUpdatePost(ctx, tx, post)
	if err != nil {
		return nil, nil, err
	}

	err = dbDeletePostBlock(ctx, tx, post.ID, mutateRequest.Removed)
	if err != nil {
		return nil, nil, err
	}

	err = dbUpdatePostBlocks(ctx, tx, mutateRequest.Edited)
	if err != nil {
		return nil, nil, err
	}

	addedIDs, err := dbCreatePostBlocks(ctx, tx, post.ID, mutateRequest.Added)
	if err != nil {
		return nil, nil, err
	}

	return res, addedIDs, nil
}

func (s *Service) DeletePost(ctx context.Context, postIDV *idv.IDV) error {
	db := s.db.WithContext(ctx)

	err := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		err := dbDeletePost(ctx, tx, postIDV)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
