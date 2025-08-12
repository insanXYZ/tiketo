package repository

import (
	"context"
	"tiketo/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] struct {
}

func (r *Repository[T]) Create(ctx context.Context, db *gorm.DB, model T) error {
	err := db.WithContext(ctx).Create(model).Error
	if err != nil {
		logger.WarnMethod("Repository.Take", err)
	}
	return err
}

func (r *Repository[T]) Take(ctx context.Context, db *gorm.DB, model T) error {
	err := db.WithContext(ctx).Where(model).Take(model).Error
	if err != nil {
		logger.WarnMethod("Repository.Take", err)
	}
	return err
}

func (r *Repository[T]) TakeForUpdate(ctx context.Context, db *gorm.DB, model T) error {
	err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where(model).Take(model).Error
	if err != nil {
		logger.WarnMethod("Repository.TakeForUpdate", err)
	}
	return err
}

func (r *Repository[T]) Delete(ctx context.Context, db *gorm.DB, model T) error {
	err := db.WithContext(ctx).Where(model).Delete(model).Error
	if err != nil {
		logger.WarnMethod("Repository.Delete", err)
	}
	return err
}

func (r *Repository[T]) Save(ctx context.Context, db *gorm.DB, model T) error {
	err := db.WithContext(ctx).Save(model).Error
	if err != nil {
		logger.WarnMethod("Repository.Save", err)
	}
	return err
}

func (r *Repository[T]) Update(ctx context.Context, db *gorm.DB, model T) error {
	err := db.WithContext(ctx).Model(model).Updates(model).Error
	if err != nil {
		logger.Warn(nil, "Error update from database :", err.Error())
	}
	return err
}
