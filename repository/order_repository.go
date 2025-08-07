package repository

import "tiketo/entity"

type OrderRepository struct {
	Repository[*entity.Order]
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}
