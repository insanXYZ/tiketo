package service

import (
	"context"
	"os"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepository       *repository.OrderRepository
	ticketRepository      *repository.TicketRepository
	orderDetailRepository *repository.OrderDetailRepository
	userRepository        *repository.UserRepository
	redis                 *redis.Client
	db                    *gorm.DB
}

func NewOrderService(orderRepository *repository.OrderRepository, orderDetailRepository *repository.OrderDetailRepository, userRepository *repository.UserRepository, ticketRepository *repository.TicketRepository, redis *redis.Client, db *gorm.DB) *OrderService {
	return &OrderService{
		orderDetailRepository: orderDetailRepository,
		userRepository:        userRepository,
		orderRepository:       orderRepository,
		ticketRepository:      ticketRepository,
		redis:                 redis,
		db:                    db,
	}
}

func (o *OrderService) HandleGetHistoryOrder(ctx context.Context, claims jwt.MapClaims, req *dto.GetOrder) (*entity.Order, error) {
	order := &entity.Order{
		ID:     req.TicketID,
		UserID: claims["sub"].(string),
	}

	err := o.orderRepository.TakeWithDetailOrder(ctx, o.db, order)
	return order, err
}

func (o *OrderService) HandleCreate(ctx context.Context, claims jwt.MapClaims, req *dto.CreateOrder) (*snap.Response, error) {
	var snapResponse *snap.Response

	err := util.ValidateStruct(req)
	if err != nil {
		return nil, err
	}

	err = o.db.Transaction(func(tx *gorm.DB) error {
		user := &entity.User{
			ID: claims["sub"].(string),
		}
		err = o.userRepository.Take(ctx, tx, user)
		if err != nil {
			return err
		}

		ticket := &entity.Ticket{
			ID: req.TicketID,
		}
		err = o.ticketRepository.Take(ctx, tx, ticket)
		if err != nil {
			return err
		}

		total := req.Quantity * ticket.Price
		order := &entity.Order{
			ID:     uuid.NewString(),
			Status: entity.Unpaid,
			UserID: user.ID,
			Total:  total,
		}

		err = o.orderRepository.Create(ctx, tx, order)
		if err != nil {
			return err
		}

		orderDetail := &entity.OrderDetail{
			OrderID:  order.ID,
			TicketId: ticket.ID,
			Quantity: uint(req.Quantity),
		}

		err = o.orderDetailRepository.Create(ctx, tx, orderDetail)
		if err != nil {
			return err
		}

		snapClient := &snap.Client{}
		snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

		transactionDetail := midtrans.TransactionDetails{
			OrderID:  order.ID,
			GrossAmt: int64(total),
		}

		customerDetail := &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		}

		items := &[]midtrans.ItemDetails{
			{
				ID:    ticket.ID,
				Qty:   int32(req.Quantity),
				Name:  ticket.Name,
				Price: int64(ticket.Price),
			},
		}

		snapReq := &snap.Request{
			TransactionDetails: transactionDetail,
			CustomerDetail:     customerDetail,
			EnabledPayments:    snap.AllSnapPaymentType,
			Items:              items,
		}

		snapResponse, err = snapClient.CreateTransaction(snapReq)
		return err

	})

	return snapResponse, err
}

func (o *OrderService) HandleGetHistoryOrders(ctx context.Context, claims jwt.MapClaims) ([]entity.Order, error) {
	var orders []entity.Order
	idUser := claims["sub"].(string)

	err := o.orderRepository.FindAllHistoryUser(ctx, o.db, &orders, idUser)
	return orders, err
}
