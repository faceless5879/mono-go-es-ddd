package adapters

import (
	"context"
	"encoding/json"
	es "github.com/faceless5879/mono-go-es-ddd/internal/common/event_sourcing"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/domain/order"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type OrderRdbRepository struct {
	db *sqlx.DB
}

func NewOrderRdbRepository(db *sqlx.DB) *OrderRdbRepository {
	return &OrderRdbRepository{db: db}
}

// SaveEvents use named return values. If return value is error instead of (err error),　defer block cannot handle query result
func (r OrderRdbRepository) SaveEvents(ctx context.Context, events []es.Event) (err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}
	defer func() {
		err = r.finishTransaction(err, tx)
	}()
	query := `INSERT INTO order_events (id, order_id, event_type, data, time_stamp) VALUES ($1,$2,$3,$4,$5)`
	for _, event := range events {
		switch e := event.(type) {
		case order.OrderCreatedEvent:
			var orderItems []models.OrderItem
			for _, oi := range e.OrderItems {
				orderItems = append(orderItems, models.OrderItem{SkuID: oi.SkuID(), Quantity: oi.Quantity()})
			}
			jsonData, _ := json.Marshal(models.OrderCreatedEventData{
				UserUUID:        e.UserUUID,
				ReceiverName:    e.ReceiverName,
				OrderItems:      orderItems,
				DeliveryAddress: e.DeliveryAddress,
			})

			if _, err := tx.Exec(query, uuid.New(), e.OrderID, e.EventType(), jsonData, e.TimeStamp()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r OrderRdbRepository) LoadEvents(ctx context.Context, aggregateID string) (events []es.Event, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return []es.Event{}, errors.Wrap(err, "unable to start transaction")
	}
	defer func() {
		err = r.finishTransaction(err, tx)
	}()

	query := `SELECT * FROM order_events WHERE order_id = $1`
	var dModels []models.OrderEvent

	rows, err := tx.Query(query, aggregateID)
	if err != nil {
		return []es.Event{}, err
	}

	for rows.Next() {
		var dm models.OrderEvent
		if err := rows.Scan(&dm.ID, &dm.OrderID, &dm.EventType, &dm.Data, &dm.TimeStamp); err != nil {
			return []es.Event{}, err
		}
		dModels = append(dModels, dm)
	}

	events, err = r.convertToDomainEvents(dModels)
	if err != nil {
		return []es.Event{}, err
	} else {
		return events, nil
	}
}

func (r OrderRdbRepository) convertToDomainEvents(events []models.OrderEvent) ([]es.Event, error) {
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

// finishTransaction rollbacks transaction if error is provided.
// If err is nil transaction is committed.
//
// If the rollback fails, we are using multierr library to add error about rollback failure.
// If the commit fails, commit error is returned.
func (r OrderRdbRepository) finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}
