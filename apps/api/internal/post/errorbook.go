package post

import (
	"errors"

	"github.com/vandad1901/p3s/packages/go/apperror"
)

var (
	ErrConflictSlug = apperror.Conflict("post.validation.SlugConflict")
	ErrPostNotFount = apperror.NotFound("post.NotFound")
)

var (
	errEmptyTitle            = errors.New("post.validation.EmptyTitle")
	errEmptySlug             = errors.New("post.validation.EmptySlug")
	errValidationInvalidSlug = errors.New("post.validation.InvalidSlug")
	errInvalidContent        = errors.New("post.validation.InvalidContent")
	errInvalidMetadata       = errors.New("post.validation.InvalidMetadata")
	errValidationBadOrdering = errors.New("post.validation.BadOrdering")
)
