package order

import (
	es "github.com/faceless5879/mono-go-es-ddd/internal/common/event_sourcing"
	"time"
)

type Order struct {
	uuid       string
	uncommited []es.Event
	version    int

	userUUID        string
	orderStatus     OrderStatus
	deliveryAddress DeliveryAddress
	orderItems      []OrderItem
	createdAt       time.Time
}

func NewOrder(ID string) (*Order, error) {
	return &Order{uuid: ID}, nil
}

func (o *Order) incrementVersion() {
	o.version++
}

func (o *Order) ID() string {
	return o.uuid
}

func (o *Order) UserID() string {
	return o.userUUID
}

func (o *Order) DeliveryAddress() DeliveryAddress {
	return o.deliveryAddress
}

func (o *Order) OrderItems() []OrderItem {
	return o.orderItems
}

func (o *Order) GetUncommitedEvents() []es.Event {
	return o.uncommited
}

func (o *Order) LoadFromHistory(events []es.Event) {
	for _, e := range events {
		o.ApplyChange(e, false)
	}
}

func (o *Order) Version() int {
	return o.version
}

func (o *Order) ApplyChange(event es.Event, loadFromHistory bool) {
	o.incrementVersion()
	switch e := event.(type) {
	case OrderCreatedEvent:
		o.orderStatus = Pending
		o.userUUID = e.UserUUID
		o.uuid = e.OrderID
		o.orderItems = e.OrderItems
		o.deliveryAddress, _ = NewDeliveryAddress(e.ReceiverName, e.DeliveryAddress)
		o.createdAt = e.TimeStamp()
		if !loadFromHistory {
			o.uncommited = append(o.uncommited, e)
		}
	}
}

const OrderCreated es.EventType = "ORDER_CREATED"

type OrderCreatedEvent struct {
	es.BaseEvent
	OrderID         string
	UserUUID        string
	ReceiverName    string
	DeliveryAddress string
	OrderItems      []OrderItem
}

func (e OrderCreatedEvent) EventType() es.EventType {
	return e.BaseEvent.EventType
}

func (e OrderCreatedEvent) TimeStamp() time.Time {
	return e.BaseEvent.TimeStamp
}

func (o *Order) Init(userUUID string, orderItems []OrderItem, deliveryAddress DeliveryAddress) {
	now := time.Now()
	o.ApplyChange(OrderCreatedEvent{
		BaseEvent: es.BaseEvent{
			EventType: OrderCreated, TimeStamp: now,
		},
		OrderID:         o.uuid,
		UserUUID:        userUUID,
		ReceiverName:    deliveryAddress.ReceiverName(),
		DeliveryAddress: deliveryAddress.Address(),
		OrderItems:      orderItems,
	}, false)
}
