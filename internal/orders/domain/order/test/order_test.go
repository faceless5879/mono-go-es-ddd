package test_test

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/domain/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder_Init(t *testing.T) {
	t.Parallel()
	item, _ := order.NewOrderItem("123", 2)
	orderItems := []order.OrderItem{item}
	deliveryAddress, _ := order.NewDeliveryAddress("Test", "Test Addr")
	userUUID := "user id"
	orderID := uuid.New().String()

	newOrder := order.NewOrder(orderID)
	newOrder.Init(userUUID, orderItems, deliveryAddress)
	events := newOrder.GetUncommitedEvents()

	assert.Equal(t, newOrder.ID(), orderID)
	assert.Equal(t, newOrder.DeliveryAddress(), deliveryAddress)
	assert.Equal(t, newOrder.UserID(), userUUID)

	// Check events
	assert.Len(t, events, 1)
	if len(events) > 0 {
		e, _ := events[0].(order.OrderCreatedEvent)
		assert.Equal(t, order.OrderCreated, e.EventType())
		assert.Equal(t, orderID, e.OrderID)
		assert.Equal(t, userUUID, e.UserUUID)
		assert.Equal(t, newOrder.DeliveryAddress().ReceiverName(), e.ReceiverName)
		assert.Equal(t, newOrder.DeliveryAddress().Address(), e.DeliveryAddress)
		assert.Equal(t, orderItems, e.OrderItems)
	}
}
