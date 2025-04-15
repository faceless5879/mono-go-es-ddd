package adapters

import (
	"context"
	"encoding/json"
	es "github.com/faceless5879/mono-go-es-ddd/internal/common/event_sourcing"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/domain/order"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRdbRepository struct {
	tx *sqlx.Tx
}

func (r *OrderRdbRepository) SaveEvents(ctx context.Context, events []es.Event) error {
	query := `INSERT INTO order_events (id, order_id, event_type, data, time_stamp) VALUES ($1,$2,$3,$4,$5)`
	for _, event := range events {
		switch e := event.(type) {
		case order.OrderCreatedEvent:
			var orderItems []models.OrderItem
			for _, oi := range e.OrderItems {
				orderItems = append(orderItems, models.OrderItem{SkuID: oi.SkuID(), Quantity: oi.Quantity()})
			}
			jsonData, err := json.Marshal(models.OrderCreatedEventData{
				UserUUID:        e.UserUUID,
				ReceiverName:    e.ReceiverName,
				OrderItems:      orderItems,
				DeliveryAddress: e.DeliveryAddress,
			})

			if err != nil {
				return err
			}

			if _, err := r.tx.Exec(query, uuid.New(), e.OrderID, e.EventType(), jsonData, e.TimeStamp()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *OrderRdbRepository) LoadEvents(ctx context.Context, aggregateID string) ([]es.Event, error) {
	query := `SELECT * FROM order_events WHERE order_id = $1`
	var dModels []models.OrderEvent

	rows, err := r.tx.Query(query, aggregateID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dm models.OrderEvent
		if err := rows.Scan(&dm.ID, &dm.OrderID, &dm.EventType, &dm.Data, &dm.TimeStamp); err != nil {
			return nil, err
		}
		dModels = append(dModels, dm)
	}

	return r.convertToDomainEvents(dModels)
}

func (r *OrderRdbRepository) convertToDomainEvents(events []models.OrderEvent) ([]es.Event, error) {
	var res []es.Event
	for _, e := range events {
		if e.EventType == "ORDER_CREATED" {
			var (
				orderItems []order.OrderItem
				data       models.OrderCreatedEventData
			)
			if err := json.Unmarshal(e.Data, &data); err != nil {
				return nil, err
			}
			for _, oi := range data.OrderItems {
				item, _ := order.NewOrderItem(oi.SkuID, oi.Quantity)
				orderItems = append(orderItems, item)
			}
			res = append(res, order.OrderCreatedEvent{
				BaseEvent: es.BaseEvent{
					EventType: order.OrderCreated,
					TimeStamp: e.TimeStamp,
				},
				OrderID:         e.OrderID,
				UserUUID:        data.UserUUID,
				ReceiverName:    data.ReceiverName,
				DeliveryAddress: data.DeliveryAddress,
				OrderItems:      orderItems,
			})
		}
	}
	return res, nil
}
