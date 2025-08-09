package repository

import "tiketo/entity"

type OrderDetailRepository struct {
	Repository[*entity.OrderDetail]
}

func NewOrderDetailRepository() *OrderDetailRepository {
	return &OrderDetailRepository{}
}
