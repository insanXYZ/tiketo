package converter

import (
	"tiketo/dto"
	"tiketo/entity"
)

func OrderEntityToDto(entity *entity.Order) *dto.Order {
	if entity == nil {
		return nil
	}

	return &dto.Order{
		ID:          entity.ID,
		Status:      string(entity.Status),
		Total:       entity.Total,
		CreatedAt:   entity.CreatedAt,
		PaidAt:      entity.PaidAt,
		OrderDetail: OrderDetailEntityToDto(entity.OrderDetail),
	}
}

func OrderEntitiesToDto(entites []entity.Order) []dto.Order {
	if entites == nil {
		return nil
	}

	orders := make([]dto.Order, len(entites))

	for i := range len(entites) {
		orders[i] = *OrderEntityToDto(&entites[i])
	}

	return orders
}
