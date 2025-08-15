package repository

import "tiketo/entity"

type OrderDetailRepositoryInterface interface {
	RepositoryInterface[*entity.OrderDetail]
}

type OrderDetailRepository struct {
	Repository[*entity.OrderDetail]
}

func NewOrderDetailRepository() *OrderDetailRepository {
	return &OrderDetailRepository{}
}
