package repository

import (
	"context"
	"tiketo/entity"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Repository[*entity.Order]
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (o *OrderRepository) FindAllHistoryUser(ctx context.Context, db *gorm.DB, dst *[]entity.Order, id string) error {
	return db.WithContext(ctx).Where("user_id = ?", id).Joins("OrderDetail").Find(dst).Error
}

func (o *OrderRepository) TakeWithDetailOrder(ctx context.Context, db *gorm.DB, dst *entity.Order) error {
	return db.WithContext(ctx).Joins("OrderDetail").Take(dst).Error
}
