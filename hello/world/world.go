package main

import (
	"encoding/json"
	"fmt"
	"github.com/faceless5879/mono-go-es-ddd/internal/common/logs"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/domain/order"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/models"
	"github.com/google/uuid"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func createNewOrder(db *sqlx.DB) {
	item, _ := order.NewOrderItem("123", 2)
	orderItems := []order.OrderItem{item}
	deliveryAddress, _ := order.NewDeliveryAddress("Test", "Test Addr")
	userUUID := "user id"
	orderID := uuid.New()

	newOrder, _ := order.NewOrder(orderID)
	newOrder.Init(userUUID, orderItems, deliveryAddress)
	events := newOrder.GetUncommitedEvents()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error, %v", err)
		return
	}
	for _, e := range events {
		switch e := e.(type) {
		case order.OrderCreatedEvent:
			query := `INSERT INTO order_events (id, order_id, event_type, data, time_stamp) VALUES ($1,$2,$3,$4,$5)`
			jsonData, _ := json.Marshal(models.EventData{
				UserUUID:        e.UserUUID,
				ReceiverName:    e.ReceiverName,
				OrderItems:      []models.OrderItem{{SkuID: e.OrderItems[0].SkuID(), Quantity: e.OrderItems[0].Quantity()}},
				DeliveryAddress: e.DeliveryAddress,
			})
			res, err := tx.Exec(query, uuid.New(), e.OrderID, e.EventType(), jsonData, e.TimeStamp())
			if err != nil {
				tx.Rollback()
				log.Fatalf("Rollback Error, %v", err)
				return
			}
			data, _ := json.Marshal(&res)
			fmt.Println("Res, ", string(data))
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error, %v", err)
	}
}

func main() {
	logs.Init()
	//app, cleanup := service.NewApplication(ctx)
	//defer cleanup()
	//
	//servers.RunHTTPServer(func(router chi.Router) http.Handler {
	//	return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	//})
	connStr := "host=localhost port=5432 user=admin password=password dbname=pg sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var events []models.OrderEvents
	query := `SELECT * FROM order_events WHERE order_id = $1`
	rows, _ := db.Query(query, "c7177dcc-88fe-4c42-bc6c-07125554889a")

	for rows.Next() {
		var (
			event    models.OrderEvents
			jsonData []byte
			data     models.EventData
		)

		if err := rows.Scan(&event.ID, &event.OrderID, &event.EventType, &jsonData, &event.TimeStamp); err != nil {
			log.Fatalf("Scan Error, %v", err)
		}
		if err := json.Unmarshal(jsonData, &data); err != nil {
			log.Fatalf("Unmarshal Error, %v", err)
		}
		event.Data = data
		events = append(events, event)
	}

	for _, e := range events {
		log.Printf("event id %v", e.ID)
		log.Printf("event type %v", e.EventType)
		log.Printf("event order id %v", e.OrderID)
		log.Printf("event data %v", e.Data)
		log.Printf("event time stamp %v", e.TimeStamp)
	}
}
