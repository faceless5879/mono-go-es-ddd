package command

import (
	"context"
	"fmt"
	"github.com/faceless5879/mono-go-es-ddd/internal/common/decorator"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/domain/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrderItem struct {
	Quantity int
	SkuID    string
}
type CreateOrder struct {
	UserUUID        string
	ReceiverName    string
	DeliveryAddress string
	OrderItems      []OrderItem
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder]

type createOrderHandler struct {
	orderRepo order.Repository
}

func NewCreateOrderHandler(orderRepo order.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder](createOrderHandler{orderRepo: orderRepo}, logger, metricsClient)
}

func (c createOrderHandler) Handle(ctx context.Context, command CreateOrder) error {
	if len(command.OrderItems) < 1 {
		return fmt.Errorf("order must have at least one item")
	}

	var orderItems []order.OrderItem
	for _, item := range command.OrderItems {
		orderItem, err := order.NewOrderItem(item.SkuID, item.Quantity)
		if err != nil {
			return err
		}
		orderItems = append(orderItems, orderItem)
	}
	deliveryAddress, err := order.NewDeliveryAddress(command.ReceiverName, command.DeliveryAddress)
	if err != nil {
		return err
	}
	newOrder := order.NewOrder(uuid.New().String())
	newOrder.Init(command.UserUUID, orderItems, deliveryAddress)
	if err := c.orderRepo.SaveEvents(ctx, newOrder.GetUncommitedEvents()); err != nil {
		return err
	}
	// TODO: publish event to MQ, run projector to create read view
	return nil
}
