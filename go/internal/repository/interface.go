package repository

import (
	"context"
	"momo-backend-go/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, c *model.Comment) error
	GetByID(ctx context.Context, id int64) (*model.Comment, error)
	GetByPostSlug(ctx context.Context, slug string) ([]*model.Comment, error)
	List(ctx context.Context, offset, limit int) ([]*model.Comment, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error
}
