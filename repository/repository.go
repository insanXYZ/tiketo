package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] struct {
}

func (r *Repository[T]) Create(ctx context.Context, db *gorm.DB, model T) error {
	return db.WithContext(ctx).Create(model).Error
}

func (r *Repository[T]) Take(ctx context.Context, db *gorm.DB, model T) error {
	return db.WithContext(ctx).Take(model).Error
}
