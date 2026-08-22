package post

import (
	"context"
	"encoding/json"

	"github.com/vandad1901/p3s/packages/go/mutatelist"
	"github.com/vandad1901/p3s/packages/go/usercontext"
	"gorm.io/gorm"
)

func validatePost(ctx context.Context, db *gorm.DB, post *Post, postBlocks []*PostBlock) error {
	err := validatePostFields(ctx, db, post)
	if err != nil {
		return err
	}

	err = validatePostBlocks(postBlocks)
	if err != nil {
		return err
	}

	return nil
}

func validatePostFields(ctx context.Context, db *gorm.DB, post *Post) error {
	err := validateTitle(post.Title)
	if err != nil {
		return err
	}

	err = validateSlug(ctx, db, post)
	if err != nil {
		return err
	}

	return nil
}

func validateTitle(title string) error {
	if len(title) == 0 {
		return errEmptyTitle
	}

	return nil
}

func validateSlug(ctx context.Context, db *gorm.DB, post *Post) error {
	err := validateSlugCharacters(post.Slug)
	if err != nil {
		return err
	}

	err = validateSlugUniqueness(ctx, db, post.ID, post.Slug)
	if err != nil {
		return err
	}

	return nil
}

func validateSlugCharacters(slug string) error {
	if len(slug) == 0 {
		return errEmptySlug
	}

	for _, r := range slug {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return errValidationInvalidSlug
		}
	}

	return nil
}

func validateSlugUniqueness(ctx context.Context, db *gorm.DB, postID int64, slug string) error {
	var exists bool

	createdBy, err := usercontext.CtxUser(ctx)
	if err != nil {
		return err
	}

	q := db.Table("post").
		Where("id <> ?", postID).
		Where("created_by = ?", createdBy).
		Where("slug = ?", slug).
		Select("1")

	err = db.Raw("SELECT EXISTS (?)", q).Scan(&exists).Error
	if err != nil {
		return err
	}

	if exists {
		return ErrConflictSlug
	}

	return nil
}

func validatePostBlocks(postBlocks []*PostBlock) error {
	for i, item := range postBlocks {
		if i != 0 {
			if postBlocks[i-1].Position >= item.Position {
				return errValidationBadOrdering
			}
		}

		err := validatePostBlockFields(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePostBlockFields(item *PostBlock) error {
	ok := json.Valid([]byte(item.Metadata))
	if !ok {
		return errInvalidMetadata
	}

	switch item.BlockType {
	case BlockTypeText:
		if item.TextContent == "" || item.MediaContent != "" {
			return errInvalidContent
		}
	case BlockTypeMedia:
		if item.MediaContent == "" || item.TextContent != "" {
			return errInvalidContent
		}
	case BlockTypeUnspecified:
		return errInvalidContent
	}

	return nil
}

func validatePostForUpdate(ctx context.Context, db *gorm.DB, post *Post, mutateList *PostBlockMutateRequest) error {
	currentPostBlocks, err := dbGetPostBlocks(ctx, db, post.ID)
	if err != nil {
		return err
	}

	finalPostBlocks, err := mutatelist.MergeItems(currentPostBlocks, mutateList)
	if err != nil {
		return err
	}

	err = validatePost(ctx, db, post, finalPostBlocks)
	if err != nil {
		return err
	}

	return nil
}
