package converter

import (
	"tiketo/dto"
	"tiketo/entity"
)

func OrderDetailEntityToDto(entity *entity.OrderDetail) *dto.OrderDetail {
	if entity == nil {
		return nil
	}

	return &dto.OrderDetail{
		Quantity: entity.Quantity,
		Ticket:   TicketEntityToDto(entity.Ticket),
	}
}
