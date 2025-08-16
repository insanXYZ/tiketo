package service

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"os"
	"slices"
	"tiketo/db"
	"tiketo/dto"
	"tiketo/dto/message"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"
	"tiketo/util/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	HandleAfterPayment(context.Context, *dto.AfterPayment) error
	HandleGetHistoryOrder(context.Context, jwt.MapClaims, *dto.GetOrder) (*entity.Order, error)
	HandleGetHistoryOrders(context.Context, jwt.MapClaims) ([]entity.Order, error)
	HandleCreate(context.Context, jwt.MapClaims, *dto.CreateOrder) (*snap.Response, error)
}

type OrderService struct {
	orderRepository       repository.OrderRepositoryInterface
	ticketRepository      repository.TicketRepositoryInterface
	orderDetailRepository repository.OrderDetailRepositoryInterface
	userRepository        repository.UserRepositoryInterface
	redis                 db.RedisInterface
	db                    *gorm.DB
}

func NewOrderService(orderRepository repository.OrderRepositoryInterface, orderDetailRepository repository.OrderDetailRepositoryInterface, userRepository repository.UserRepositoryInterface, ticketRepository repository.TicketRepositoryInterface, redis db.RedisInterface, db *gorm.DB) *OrderService {
	return &OrderService{
		orderDetailRepository: orderDetailRepository,
		userRepository:        userRepository,
		orderRepository:       orderRepository,
		ticketRepository:      ticketRepository,
		redis:                 redis,
		db:                    db,
	}
}

func (o *OrderService) HandleAfterPayment(ctx context.Context, req *dto.AfterPayment) error {
	logger.EnteringMethod("OrderService.HandleAfterPayment")
	err := util.ValidateStruct(req)
	if err != nil {
		return err
	}

	logger.Debug(logrus.Fields{
		"req": req,
	}, "Success retrieve request on OrderService.HandleAfterPayment")

	validTransactionStatus := []string{"capture", "settlement"}

	if !slices.Contains(validTransactionStatus, req.TransactionStatus) {
		return errors.New(message.ErrTransactionStatus)
	}

	s := sha512.New()
	s.Write([]byte(req.OrderId + req.StatusCode + req.GrossAmount + os.Getenv("MIDTRANS_SERVER_KEY")))

	byteSignatureKey := hex.EncodeToString(s.Sum(nil))

	if string(byteSignatureKey) != req.SignatureKey {
		return errors.New(message.ErrSignatureKey)
	}

	err = o.orderRepository.Transaction(ctx, o.db, func(tx *gorm.DB) error {
		order := &entity.Order{
			ID: req.OrderId,
		}

		err = o.orderRepository.TakeWithDetailOrder(ctx, tx, order)
		if err != nil {
			return err
		}

		now := time.Now()

		order.Status = entity.Paid
		order.PaidAt = &now
		err = o.orderRepository.Save(ctx, tx, order)
		if err != nil {
			return err
		}

		ticket := &entity.Ticket{
			ID: order.OrderDetail.TicketId,
		}

		err = o.ticketRepository.Take(ctx, tx, ticket)
		if err != nil {
			return err
		}

		ticket.Quantity = ticket.Quantity - int(order.OrderDetail.Quantity)

		return o.ticketRepository.Save(ctx, tx, ticket)
	})

	return err
}

func (o *OrderService) HandleGetHistoryOrder(ctx context.Context, claims jwt.MapClaims, req *dto.GetOrder) (*entity.Order, error) {
	logger.EnteringMethod("OrderService.HandleGetHistoryOrder")
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

	err = o.orderRepository.Transaction(ctx, o.db, func(tx *gorm.DB) error {
		user := &entity.User{
			ID: claims["sub"].(string),
		}
		err := o.userRepository.Take(ctx, tx, user)
		if err != nil {
			return err
		}

		ticket := &entity.Ticket{
			ID: req.TicketID,
		}

		err = o.ticketRepository.TakeForUpdate(ctx, tx, ticket)
		if err != nil {
			return err
		}

		if ticket.Quantity < req.Quantity {
			return errors.New(message.ErrQuantityOrder)
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

		res, errSnap := snapClient.CreateTransaction(snapReq)
		if errSnap != nil {
			return errSnap
		}

		snapResponse = res
		return nil
	})

	return snapResponse, err
}

func (o *OrderService) HandleGetHistoryOrders(ctx context.Context, claims jwt.MapClaims) ([]entity.Order, error) {
	var orders []entity.Order
	idUser := claims["sub"].(string)

	err := o.orderRepository.FindAllOrderHistoryUser(ctx, o.db, &orders, idUser)
	return orders, err
}
