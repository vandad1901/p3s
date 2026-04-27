package post

import (
	"purpl3shadow/gen/commonpb"

	"gorm.io/gorm"
)

func dbCreate(tx *gorm.DB, post *Post) (*commonpb.IDVersion, error) {
	err := tx.Create(post).Error
	if err != nil {
		return nil, err
	}

	return &commonpb.IDVersion{Id: post.ID}, nil
}

func dbGetByID(tx *gorm.DB, id int64) (*Post, error) {
	post := new(Post)

	err := tx.Take(post, id).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func dbUpdate(tx *gorm.DB, post *Post) (*Post, error) {
	err := tx.Where("id = ?", post.ID).
		Updates(map[string]any{
			"content": post.Content,
		}).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func dbDelete(tx *gorm.DB, id int64) error {
	err := tx.Delete(&Post{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
