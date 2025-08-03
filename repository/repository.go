package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func (r *Repository[T]) Create(ctx context.Context, model T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *Repository[T]) Take(ctx context.Context, model T) error {
	return r.db.WithContext(ctx).Take(model).Error
}
