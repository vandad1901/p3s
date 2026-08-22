package post

import (
	"context"
	"errors"
	"time"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/idv"
	"github.com/vandad1901/p3s/packages/go/usercontext"
	"gorm.io/gorm"
)

func dbCreatePost(ctx context.Context, db *gorm.DB, p *Post) (*idv.IDV, error) {
	currentTime := time.Now()

	currentUser, err := usercontext.CtxUser(ctx)
	if err != nil {
		return nil, err
	}

	p.ID = 0
	p.CreatedAt = currentTime
	p.CreatedBy = currentUser
	p.UpdatedAt = currentTime
	p.UpdatedBy = currentUser

	err = db.Create(p).Error
	if err != nil {
		if dbpattern.IsConstraintViolation(err, "post_slug_key") {
			return nil, ErrConflictSlug
		}

		return nil, err
	}

	return &idv.IDV{
		ID:        p.ID,
		UpdatedAt: currentTime,
	}, nil
}

func dbCreatePostBlocks(_ context.Context, db *gorm.DB, postID int64, items []*PostBlock) ([]int64, error) {
	createdIDs := make([]int64, len(items))
	if len(items) == 0 {
		return createdIDs, nil
	}

	for _, item := range items {
		item.PostID = postID
	}

	err := db.Create(items).Error
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		createdIDs[i] = item.ID
	}

	return createdIDs, nil
}

func dbGetPost(_ context.Context, db *gorm.DB, postID int64) (*Post, error) {
	res := new(Post)

	err := db.Table("post").
		Where("id = ?", postID).
		First(res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFount
		}

		return nil, err
	}

	return res, nil
}

func dbGetPostBlocks(_ context.Context, db *gorm.DB, postID int64) ([]*PostBlock, error) {
	var res []*PostBlock

	err := db.Table("post_block").
		Where("post_id = ?", postID).
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func dbUpdatePost(ctx context.Context, db *gorm.DB, p *Post) (*idv.IDV, error) {
	currentTime := time.Now()

	currentUser, err := usercontext.CtxUser(ctx)
	if err != nil {
		return nil, err
	}

	res := db.Model(Post{}).
		Where("id = ?", p.ID).
		Where("updated_at = ?", p.UpdatedAt).
		Updates(
			map[string]any{
				"title":  p.Title,
				"slug":   p.Slug,
				"status": p.Status,

				"updated_at": currentTime,
				"updated_by": currentUser,
			},
		)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected != 1 {
		return nil, ErrPostNotFount
	}

	return &idv.IDV{
		ID:        p.ID,
		UpdatedAt: currentTime,
	}, nil
}

func dbUpdatePostBlocks(_ context.Context, db *gorm.DB, items []*PostBlock) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		res := db.Model(PostBlock{}).
			Where("id = ?", item.ID).
			Where("post_id = ?", item.PostID).
			Where("block_type = ?", BlockTypeText). // maybe just Text-like blocks?
			Updates(map[string]any{
				"position": item.Position,

				"text_content":  item.TextContent,
				"media_content": item.MediaContent,
				"metadata":      item.Metadata,
			})
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return ErrPostNotFount
		}
	}

	return nil
}

func dbDeletePostBlock(_ context.Context, db *gorm.DB, postID int64, items []int64) error {
	if len(items) == 0 {
		return nil
	}

	err := db.
		Where("post_id = ?", postID).
		Where("id in (?)", items).
		Delete(PostBlock{}).Error
	if err != nil {
		return err
	}

	return nil
}

func dbDeletePost(_ context.Context, db *gorm.DB, postIDV *idv.IDV) error {
	res := db.
		Where("id = ?", postIDV.ID).
		Where("updated_at = ?", postIDV.UpdatedAt).
		Delete(Post{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 1 {
		return ErrPostNotFount
	}

	return nil
}
