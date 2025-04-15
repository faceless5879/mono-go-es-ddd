package models

import (
	"encoding/json"
	"time"
)

type OrderItem struct {
	SkuID    string `json:"sku_id"`
	Quantity int    `json:"quantity"`
}
type OrderCreatedEventData struct {
	UserUUID        string      `json:"user_uuid"`
	ReceiverName    string      `json:"receiver_name"`
	OrderItems      []OrderItem `json:"order_items"`
	DeliveryAddress string      `json:"delivery_address"`
}

type OrderEvent struct {
	ID        string          `db:"id"`
	OrderID   string          `db:"order_id"`
	EventType string          `db:"event_type"`
	Data      json.RawMessage `json:"data"`
	TimeStamp time.Time       `db:"time_stamp"`
}
