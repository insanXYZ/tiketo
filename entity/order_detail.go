package entity

type OrderDetail struct {
	ID       uint    `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID  string  `gorm:"column:order_id"`
	TicketId string  `gorm:"column:ticket_id"`
	Quantity uint    `gorm:"column:quantity"`
	Order    *Order  `gorm:"foreignKey:order_id;references:id"`
	Ticket   *Ticket `gorm:"foreignKey:ticket_id;references:id"`
}

func (o OrderDetail) TableName() string {
	return "order_details"
}
